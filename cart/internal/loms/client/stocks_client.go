package client

import (
	"context"

	loms "route256/cart/pb/api"
)

type StocksClient struct {
	client loms.StocksClient
}

func NewStocksClient(client loms.StocksClient) *StocksClient {
	return &StocksClient{client: client}
}

func (sc *StocksClient) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	return sc.client.StocksInfo(ctx, req)
}
