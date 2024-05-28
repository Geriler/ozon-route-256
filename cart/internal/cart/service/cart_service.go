package service

import (
	"route256/cart/internal/cart/model"
)

type CartRepository interface {
	AddItems(userID model.UserID, item model.Item)
	DeleteItems(userID model.UserID, itemID model.SkuID)
	DeleteCart(userID model.UserID)
	GetCart(userID model.UserID) *model.Cart
}

type CartService struct {
	cartRepository CartRepository
}

func NewCartService(cartRepository CartRepository) *CartService {
	return &CartService{
		cartRepository: cartRepository,
	}
}

func (cs *CartService) GetCartByUserID(userID model.UserID) *model.Cart {
	return cs.cartRepository.GetCart(userID)
}

func (cs *CartService) AddItemsToCart(userID model.UserID, item model.Item) {
	cs.cartRepository.AddItems(userID, item)
}

func (cs *CartService) DeleteItemsFromCart(userID model.UserID, itemID model.SkuID) {
	cs.cartRepository.DeleteItems(userID, itemID)
}

func (cs *CartService) DeleteCartByUserID(userID model.UserID) {
	cs.cartRepository.DeleteCart(userID)
}

func (cs *CartService) GetTotalPrice(cart *model.Cart) uint32 {
	totalPrice := uint32(0)
	for _, item := range cart.Items {
		totalPrice += item.Price * uint32(item.Count)
	}
	return totalPrice
}
