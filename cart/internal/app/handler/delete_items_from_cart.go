package handler

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteItemsFromCart(ctx context.Context, req *model.UserSKURequest) error {
	ctx, span := h.tracer.Start(ctx, "DeleteItemsFromCart", trace.WithAttributes(
		attribute.Int("user_id", int(req.UserID)),
		attribute.Int("sku_id", int(req.SKU)),
	))
	defer span.End()

	span.AddEvent("Delete items from cart by user id and SKU")
	h.cartService.DeleteItemsFromCart(ctx, req.UserID, req.SKU)

	span.AddEvent("Return success message")
	return nil
}
