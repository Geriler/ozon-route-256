package app

import (
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256/cart/internal/config"
	"route256/cart/internal/loms/client"
	loms "route256/loms/pb/api"
)

func NewGRPCClient(cfg config.Config) (*client.GRPCClient, error) {
	conn, err := grpc.NewClient(":"+strconv.Itoa(cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &client.GRPCClient{}, err
	}

	oc := loms.NewOrderClient(conn)
	orderClient := client.NewOrderClient(oc)

	sc := loms.NewStocksClient(conn)
	stocksClient := client.NewStocksClient(sc)

	return client.NewGRPCClient(orderClient, stocksClient), nil
}
