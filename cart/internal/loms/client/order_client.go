package client

import (
	"context"
	"time"

	loms "route256/cart/pb/api"
)

type OrderClient struct {
	client  loms.OrderClient
	timeout time.Duration
}

func NewOrderClient(client loms.OrderClient, timeout time.Duration) *OrderClient {
	return &OrderClient{client: client, timeout: timeout}
}

func (oc *OrderClient) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, oc.timeout)
	defer cancel()

	return oc.client.OrderCreate(ctx, req)
}

func (oc *OrderClient) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, oc.timeout)
	defer cancel()

	return oc.client.OrderInfo(ctx, req)
}

func (oc *OrderClient) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, oc.timeout)
	defer cancel()

	return oc.client.OrderPay(ctx, req)
}

func (oc *OrderClient) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, oc.timeout)
	defer cancel()

	return oc.client.OrderCancel(ctx, req)
}
