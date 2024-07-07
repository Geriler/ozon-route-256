package handler

import (
	"context"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteCart(ctx context.Context, req *model.UserRequest) error {
	ctx, span := h.tracer.Start(ctx, "DeleteCart")
	defer span.End()

	span.AddEvent("Delete cart by user id")
	h.cartService.DeleteCartByUserID(ctx, req.UserID)

	span.AddEvent("Return success message")
	return nil
}
