package handler

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/cart/internal/cart/model"
	loms "route256/cart/pb/api"
	"route256/cart/pkg/lib/tracing"
)

var ErrNotEnoughStock = errors.New("not enough stock")

func (h *CartHandler) AddItemsToCart(ctx context.Context, req *model.UserSKUCountRequest) error {
	ctx, span := tracing.StartSpanFromContext(ctx, "AddItemsToCart", trace.WithAttributes(
		attribute.Int("user_id", int(req.UserID)),
		attribute.Int("sku_id", int(req.SKU)),
		attribute.Int("count", int(req.Count)),
	))
	defer span.End()
	ctx = tracing.InjectSpanContext(ctx, span.SpanContext())

	span.AddEvent("Get product from ProductService")
	product, err := h.productService.GetProduct(ctx, req.SKU)
	if err != nil {
		return err
	}

	span.AddEvent("Get stocks from Loms")
	stocksInfo, err := h.grpcClient.Stocks.StocksInfo(ctx, &loms.StocksInfoRequest{SkuId: int64(req.SKU)})
	if err != nil {
		return err
	}

	span.AddEvent("Check if there is enough stock")
	cart, err := h.cartService.GetCartByUserID(ctx, req.UserID)
	if stocksInfo.GetCount() < req.Count || err == nil && cart.Items[req.SKU].Count+req.Count > stocksInfo.GetCount() {
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
