package handler

import (
	"route256/cart/internal/cart/model"
	productModel "route256/cart/internal/product/model"
)

type CartService interface {
	AddItemsToCart(userID model.UserID, item model.Item)
	DeleteItemsFromCart(userID model.UserID, itemID model.SkuID)
	DeleteCartByUserID(userID model.UserID)
	GetCartByUserID(userID model.UserID) (*model.Cart, error)
	GetTotalPrice(cart *model.Cart) uint32
}

type ProductService interface {
	GetProduct(skuId model.SkuID) (*productModel.Product, error)
}

type CartHandler struct {
	cartService    CartService
	productService ProductService
}

func NewCartHandler(cartService CartService, productService ProductService) *CartHandler {
	return &CartHandler{
		cartService:    cartService,
		productService: productService,
	}
}
