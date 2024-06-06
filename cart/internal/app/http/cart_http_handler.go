package http

import (
	"log/slog"

	"route256/cart/internal/cart/model"
)

type CartHandler interface {
	AddItemsToCart(req *model.UserSKUCountRequest) error
	DeleteItemsFromCart(req *model.UserSKURequest) error
	DeleteCart(req *model.UserRequest) error
	GetCart(req *model.UserRequest) (model.CartResponse, error)
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
