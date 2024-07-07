package handler

import (
	"context"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteItemsFromCart(ctx context.Context, req *model.UserSKURequest) error {
	ctx, span := h.tracer.Start(ctx, "DeleteItemsFromCart")
	defer span.End()

	span.AddEvent("Delete items from cart by user id and SKU")
	h.cartService.DeleteItemsFromCart(ctx, req.UserID, req.SKU)

	span.AddEvent("Return success message")
	return nil
}
