package client

import (
	"context"
	"time"

	loms "route256/cart/pb/api"
)

type StocksClient struct {
	client  loms.StocksClient
	timeout time.Duration
}

func NewStocksClient(client loms.StocksClient, timeout time.Duration) *StocksClient {
	return &StocksClient{client: client, timeout: timeout}
}

func (sc *StocksClient) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, sc.timeout)
	defer cancel()

	return sc.client.StocksInfo(ctx, req)
}
