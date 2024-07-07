package handler

import (
	"context"
	"errors"

	"route256/cart/internal/cart/model"
	loms "route256/cart/pb/api"
)

var ErrNotEnoughStock = errors.New("not enough stock")

func (h *CartHandler) AddItemsToCart(ctx context.Context, req *model.UserSKUCountRequest) error {
	ctx, span := h.tracer.Start(ctx, "AddItemsToCart")
	defer span.End()

	span.AddEvent("Get product from ProductService")
	product, err := h.productService.GetProduct(req.SKU)
	if err != nil {
		return err
	}

	span.AddEvent("Get stocks from Loms")
	stocksInfo, err := h.grpcClient.Stocks.StocksInfo(ctx, &loms.StocksInfoRequest{SkuId: int64(req.SKU)})
	if err != nil {
		return err
	}

	if stocksInfo.GetCount() < req.Count {
		return ErrNotEnoughStock
	}

	item := model.Item{
		SKU:   req.SKU,
		Name:  product.Name,
		Count: req.Count,
		Price: product.Price,
	}

	span.AddEvent("Add items to cart")
	h.cartService.AddItemsToCart(ctx, req.UserID, item)

	span.AddEvent("Return success message")
	return nil
}
