package repository

import (
	"route256/cart/internal/cart/model"
)

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
	cart := r.getCart(userID)

	_, itemExists := cart.Items[item.SKU]
	if !itemExists {
		items := &model.Item{
			SKU:   item.SKU,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		}
		cart.Items[item.SKU] = items

		return
	}

	cart.Items[item.SKU].Count += item.Count

	return
}

func (r *InMemoryCartRepository) DeleteItems(userID model.UserID, itemID model.SkuID) {
	cart := r.getCart(userID)

	_, itemExists := cart.Items[itemID]
	if !itemExists {
		return
	} else {
		delete(cart.Items, itemID)
	}

	return
}

func (r *InMemoryCartRepository) DeleteCart(userID model.UserID) {
	cart := &model.Cart{
		Items: make(map[model.SkuID]*model.Item),
	}

	r.carts[userID] = cart

	return
}

func (r *InMemoryCartRepository) GetCart(userID model.UserID) *model.Cart {
	cart := r.getCart(userID)
	return cart
}

func (r *InMemoryCartRepository) getCart(userID model.UserID) *model.Cart {
	cart, cartExists := r.carts[userID]
	if !cartExists {
		cart = &model.Cart{
			Items: make(map[model.SkuID]*model.Item),
		}

		r.carts[userID] = cart
	}
	return cart
}
