package client

import (
	"context"

	loms "route256/loms/pb/api"
)

type OrderClient struct {
	client loms.OrderClient
}

func NewOrderClient(client loms.OrderClient) *OrderClient {
	return &OrderClient{client: client}
}

func (oc *OrderClient) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	return oc.client.OrderCreate(ctx, req)
}

func (oc *OrderClient) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	return oc.client.OrderInfo(ctx, req)
}

func (oc *OrderClient) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	return oc.client.OrderPay(ctx, req)
}

func (oc *OrderClient) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {
	return oc.client.OrderCancel(ctx, req)
}
