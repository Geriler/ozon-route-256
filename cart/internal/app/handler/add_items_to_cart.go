package handler

import (
	"context"
	"errors"

	"route256/cart/internal/cart/model"
	loms "route256/loms/pb/api"
)

const ErrNotEnoughStock = "not enough stock"

func (h *CartHandler) AddItemsToCart(ctx context.Context, req *model.UserSKUCountRequest) error {
	product, err := h.productService.GetProduct(req.SKU)
	if err != nil {
		return err
	}

	stocksInfo, err := h.grpcClient.Stocks.StocksInfo(ctx, &loms.StocksInfoRequest{SkuId: int64(req.SKU)})
	if err != nil {
		return err
	}

	if stocksInfo.GetCount() < req.Count {
		return errors.New(ErrNotEnoughStock)
	}

	item := model.Item{
		SKU:   req.SKU,
		Name:  product.Name,
		Count: req.Count,
		Price: product.Price,
	}

	h.cartService.AddItemsToCart(ctx, req.UserID, item)

	return nil
}
