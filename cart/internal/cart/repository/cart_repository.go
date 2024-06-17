package repository

import (
	"context"
	"errors"
	"sync"

	"route256/cart/internal/cart/model"
)

var ErrCartNotFoundOrEmpty = errors.New("cart not found or empty")

type InMemoryCartRepository struct {
	carts map[model.UserID]*model.Cart
	mutex *sync.RWMutex
}

func NewInMemoryCartRepository() *InMemoryCartRepository {
	carts := make(map[model.UserID]*model.Cart)
	return &InMemoryCartRepository{
		carts: carts,
		mutex: &sync.RWMutex{},
	}
}

func (r *InMemoryCartRepository) AddItems(_ context.Context, userID model.UserID, item model.Item) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	cart, cartExists := r.carts[userID]
	if !cartExists {
		cart = &model.Cart{
			Items: make(map[model.SkuID]*model.Item),
		}

		r.carts[userID] = cart
	}

	_, itemExists := cart.Items[item.SKU]
	if itemExists {
		cart.Items[item.SKU].Count += item.Count
	} else {
		items := &model.Item{
			SKU:   item.SKU,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		}
		cart.Items[item.SKU] = items
	}
}

func (r *InMemoryCartRepository) DeleteItems(_ context.Context, userID model.UserID, itemID model.SkuID) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	cart, cartExists := r.carts[userID]
	if !cartExists {
		return
	}

	_, itemExists := cart.Items[itemID]
	if itemExists {
		delete(cart.Items, itemID)
	}
}

func (r *InMemoryCartRepository) DeleteCart(_ context.Context, userID model.UserID) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.carts, userID)
}

func (r *InMemoryCartRepository) GetCart(_ context.Context, userID model.UserID) (*model.Cart, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	cart, cartExists := r.carts[userID]
	if !cartExists {
		return nil, ErrCartNotFoundOrEmpty
	}
	if len(cart.Items) == 0 {
		return nil, ErrCartNotFoundOrEmpty
	}
	return cart, nil
}
