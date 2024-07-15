package handler

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/cart/internal/cart/model"
	"route256/cart/pkg/lib/tracing"
)

func (h *CartHandler) DeleteCart(ctx context.Context, req *model.UserRequest) error {
	ctx, span := tracing.StartSpanFromContext(ctx, "DeleteCart", trace.WithAttributes(
		attribute.Int("user_id", int(req.UserID)),
	))
	defer span.End()

	span.AddEvent("Delete cart by user id")
	h.cartService.DeleteCartByUserID(ctx, req.UserID)

	span.AddEvent("Return success message")
	return nil
}
