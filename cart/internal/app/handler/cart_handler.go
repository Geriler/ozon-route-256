package handler

import (
	"context"

	"route256/cart/internal/cart/model"
	"route256/cart/internal/loms/client"
	productModel "route256/cart/internal/product/model"
)

type CartService interface {
	AddItemsToCart(ctx context.Context, userID model.UserID, item model.Item)
	DeleteItemsFromCart(ctx context.Context, userID model.UserID, itemID model.SkuID)
	DeleteCartByUserID(ctx context.Context, userID model.UserID)
	GetCartByUserID(ctx context.Context, userID model.UserID) (*model.Cart, error)
	GetTotalPrice(ctx context.Context, cart *model.Cart) uint32
}

type ProductService interface {
	GetProduct(skuId model.SkuID) (*productModel.Product, error)
}

type CartHandler struct {
	cartService    CartService
	productService ProductService
	grpcClient     client.GRPCClient
}

func NewCartHandler(cartService CartService, productService ProductService, grpcClient client.GRPCClient) *CartHandler {
	return &CartHandler{
		cartService:    cartService,
		productService: productService,
		grpcClient:     grpcClient,
	}
}
