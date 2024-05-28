package model

import "route256/cart/internal/cart/model"

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type GetProductRequest struct {
	Token string      `json:"token"`
	Sku   model.SkuID `json:"sku"`
}

type GetProductErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
