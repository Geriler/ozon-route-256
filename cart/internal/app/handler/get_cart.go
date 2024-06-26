package handler

import (
	"context"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) GetCart(ctx context.Context, req *model.UserRequest) (model.CartResponse, error) {
	cart, err := h.cartService.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return model.CartResponse{}, err
	}

	eg, egCtx := errgroup.WithContext(ctx)
	ticker := time.NewTicker(time.Second / time.Duration(h.productService.GetRPS()))
	for _, item := range cart.Items {
		eg.Go(func() error {
			select {
			case <-egCtx.Done():
				return nil
			case <-ticker.C:
				item := item
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
