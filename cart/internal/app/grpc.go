package app

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256/cart/internal/config"
	"route256/cart/internal/loms/client"
	loms "route256/cart/pb/api"
)

func NewGRPCClient(cfg config.Config) (*client.GRPCClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("dns:%s:%d", cfg.GRPC.Host, cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	oc := loms.NewOrderClient(conn)
	orderClient := client.NewOrderClient(oc, cfg.GRPC.Timeout)

	sc := loms.NewStocksClient(conn)
	stocksClient := client.NewStocksClient(sc, cfg.GRPC.Timeout)

	return client.NewGRPCClient(orderClient, stocksClient, conn), nil
}
