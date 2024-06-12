package service

import (
	"context"

	"route256/cart/internal/cart/model"
)

type CartRepository interface {
	AddItems(ctx context.Context, userID model.UserID, item model.Item)
	DeleteItems(ctx context.Context, userID model.UserID, itemID model.SkuID)
	DeleteCart(ctx context.Context, userID model.UserID)
	GetCart(ctx context.Context, userID model.UserID) (*model.Cart, error)
}

type CartService struct {
	cartRepository CartRepository
}

func NewCartService(cartRepository CartRepository) *CartService {
	return &CartService{
		cartRepository: cartRepository,
	}
}

func (cs *CartService) GetCartByUserID(ctx context.Context, userID model.UserID) (*model.Cart, error) {
	return cs.cartRepository.GetCart(ctx, userID)
}

func (cs *CartService) AddItemsToCart(ctx context.Context, userID model.UserID, item model.Item) {
	cs.cartRepository.AddItems(ctx, userID, item)
}

func (cs *CartService) DeleteItemsFromCart(ctx context.Context, userID model.UserID, itemID model.SkuID) {
	cs.cartRepository.DeleteItems(ctx, userID, itemID)
}

func (cs *CartService) DeleteCartByUserID(ctx context.Context, userID model.UserID) {
	cs.cartRepository.DeleteCart(ctx, userID)
}

func (cs *CartService) GetTotalPrice(_ context.Context, cart *model.Cart) uint32 {
	var totalPrice uint32
	for _, item := range cart.Items {
		totalPrice += item.Price * uint32(item.Count)
	}
	return totalPrice
}
