package handler

import (
	"log/slog"

	"route256/cart/internal/cart/model"
	productModel "route256/cart/internal/product/model"
)

type CartService interface {
	AddItemsToCart(userID model.UserID, item model.Item)
	DeleteItemsFromCart(userID model.UserID, itemID model.SkuID)
	DeleteCartByUserID(userID model.UserID)
	GetCartByUserID(userID model.UserID) *model.Cart
	GetTotalPrice(cart *model.Cart) uint32
}

type ProductService interface {
	GetProduct(skuId model.SkuID) (*productModel.Product, error)
}

type CartHandler struct {
	cartService    CartService
	productService ProductService
	logger         *slog.Logger
}

func NewCartHandler(cartService CartService, productService ProductService, logger *slog.Logger) *CartHandler {
	return &CartHandler{
		cartService:    cartService,
		productService: productService,
		logger:         logger,
	}
}
