package handler

import (
	"context"
	"sort"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/cart/internal/cart/model"
	"route256/cart/pkg/lib/errgroup"
	"route256/cart/pkg/lib/tracing"
)

func (h *CartHandler) GetCart(ctx context.Context, req *model.UserRequest) (model.CartResponse, error) {
	ctx, span := tracing.StartSpanFromContext(ctx, "GetCart", trace.WithAttributes(
		attribute.Int("user_id", int(req.UserID)),
	))
	defer span.End()

	span.AddEvent("Get cart by user id")
	cart, err := h.cartService.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return model.CartResponse{}, err
	}

	span.AddEvent("Get product information from ProductService")
	eg, egCtx := errgroup.WithContext(ctx)
	ticker := time.NewTicker(time.Second / time.Duration(h.productService.GetRPSLimit()))
	for _, item := range cart.Items {
		eg.Go(func() error {
			var err error
			item := item
			select {
			case <-egCtx.Done():
				return nil
			case <-ticker.C:
				_, err = h.productService.GetProduct(ctx, item.SKU)
				if err != nil {
					egCtx.Done()
					return err
				}
				return nil
			}
		})
	}
	if err := eg.Wait(); err != nil {
		return model.CartResponse{}, err
	}

	span.AddEvent("Calculate total price")
	totalPrice := h.cartService.GetTotalPrice(ctx, cart)
	items := make([]model.Item, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})

	span.AddEvent("Return cart information")
	return model.CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
