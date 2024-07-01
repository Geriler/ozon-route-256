package handler

import (
	"context"
	"sort"
	"time"

	"route256/cart/internal/cart/model"
	"route256/cart/pkg/lib/errgroup"
)

func (h *CartHandler) GetCart(ctx context.Context, req *model.UserRequest) (model.CartResponse, error) {
	cart, err := h.cartService.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return model.CartResponse{}, err
	}

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
				_, err = h.productService.GetProduct(item.SKU)
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

	totalPrice := h.cartService.GetTotalPrice(ctx, cart)
	items := make([]model.Item, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})
	return model.CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
