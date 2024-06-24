package model

import "errors"

var (
	ErrSkuNotFound    = errors.New("sku not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)
