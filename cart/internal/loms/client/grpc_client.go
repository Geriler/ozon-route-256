package client

import (
	"context"

	"google.golang.org/grpc"
	loms "route256/cart/pb/api"
)

type GRPCClient struct {
	Order  OrderService
	Stocks StocksService
	conn   *grpc.ClientConn
}

func NewGRPCClient(order OrderService, stocks StocksService, conn *grpc.ClientConn) *GRPCClient {
	return &GRPCClient{
		Order:  order,
		Stocks: stocks,
		conn:   conn,
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

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}
