package handler

import (
	"context"
	"fmt"

	"route256/cart/internal/cart/model"
	loms "route256/cart/pb/api"
)

var ErrNotEnoughStock = fmt.Errorf("%s", "not enough stock")

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
		return ErrNotEnoughStock
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
