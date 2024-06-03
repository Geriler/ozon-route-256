package repository

import (
	"errors"

	"route256/cart/internal/cart/model"
)

const ErrCartNotFoundOrEmpty = "cart not found or empty"

type InMemoryCartRepository struct {
	carts map[model.UserID]*model.Cart
}

func NewInMemoryCartRepository() *InMemoryCartRepository {
	carts := make(map[model.UserID]*model.Cart)
	return &InMemoryCartRepository{
		carts: carts,
	}
}

func (r *InMemoryCartRepository) AddItems(userID model.UserID, item model.Item) {
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

func (r *InMemoryCartRepository) DeleteItems(userID model.UserID, itemID model.SkuID) {
	cart, cartExists := r.carts[userID]
	if !cartExists {
		return
	}

	_, itemExists := cart.Items[itemID]
	if itemExists {
		delete(cart.Items, itemID)
	}
}

func (r *InMemoryCartRepository) DeleteCart(userID model.UserID) {
	delete(r.carts, userID)
}

func (r *InMemoryCartRepository) GetCart(userID model.UserID) (*model.Cart, error) {
	cart, cartExists := r.carts[userID]
	if !cartExists {
		return nil, errors.New(ErrCartNotFoundOrEmpty)
	}
	if len(cart.Items) == 0 {
		return nil, errors.New(ErrCartNotFoundOrEmpty)
	}
	return cart, nil
}
