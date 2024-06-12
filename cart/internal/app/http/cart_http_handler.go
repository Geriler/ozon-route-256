package http

import (
	"context"
	"log/slog"

	"route256/cart/internal/cart/model"
)

type CartHandler interface {
	AddItemsToCart(ctx context.Context, req *model.UserSKUCountRequest) error
	DeleteItemsFromCart(ctx context.Context, req *model.UserSKURequest) error
	DeleteCart(ctx context.Context, req *model.UserRequest) error
	GetCart(ctx context.Context, req *model.UserRequest) (model.CartResponse, error)
	Checkout(ctx context.Context, req *model.UserRequest) (model.CartCheckoutResponse, error)
}

type CartHttpHandlers struct {
	cartHandler CartHandler
	logger      *slog.Logger
}

func NewCartHttpHandlers(cartHandler CartHandler, logger *slog.Logger) *CartHttpHandlers {
	return &CartHttpHandlers{
		cartHandler: cartHandler,
		logger:      logger,
	}
}
