package handler

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/cart/internal/cart/model"
	loms "route256/cart/pb/api"
)

func (h *CartHandler) Checkout(ctx context.Context, req *model.UserRequest) (model.CartCheckoutResponse, error) {
	ctx, span := h.tracer.Start(ctx, "Checkout", trace.WithAttributes(
		attribute.Int("user_id", int(req.UserID)),
	))
	defer span.End()

	span.AddEvent("Get cart by user id")
	cart, err := h.cartService.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return model.CartCheckoutResponse{}, err
	}

	span.AddEvent("Convert cart items to Loms items")
	items := make([]*loms.Item, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, &loms.Item{
			SkuId: int64(item.SKU),
			Count: item.Count,
		})
	}

	span.AddEvent("Create order in Loms")
	orderCreateResponse, err := h.grpcClient.Order.OrderCreate(ctx, &loms.OrderCreateRequest{
		UserId: int64(req.UserID),
		Items:  items,
	})
	if err != nil {
		return model.CartCheckoutResponse{}, err
	}

	span.AddEvent("Delete cart by user id")
	h.cartService.DeleteCartByUserID(ctx, req.UserID)

	span.AddEvent("Return order id")
	return model.CartCheckoutResponse{
		OrderID: orderCreateResponse.OrderId,
	}, nil
}
