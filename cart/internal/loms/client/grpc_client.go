package client

import (
	"context"

	loms "route256/cart/pb/api"
)

type GRPCClient struct {
	Order  OrderService
	Stocks StocksService
}

func NewGRPCClient(order OrderService, stocks StocksService) *GRPCClient {
	return &GRPCClient{
		Order:  order,
		Stocks: stocks,
	}
}

type OrderService interface {
	OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error)
	OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error)
	OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error)
	OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error)
}

type StocksService interface {
	StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error)
}
