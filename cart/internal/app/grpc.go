package app

import (
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256/cart/internal/config"
	"route256/cart/internal/loms/client"
	lomsService "route256/cart/internal/loms/service"
	loms "route256/loms/pb/api"
)

func NewGRPCClient(cfg config.Config) (*lomsService.LomsService, error) {
	conn, err := grpc.NewClient(":"+strconv.Itoa(cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &lomsService.LomsService{}, err
	}

	oc := loms.NewOrderClient(conn)
	orderClient := client.NewOrderClient(oc)

	sc := loms.NewStocksClient(conn)
	stocksClient := client.NewStocksClient(sc)

	return lomsService.NewLomsService(orderClient, stocksClient), nil
}
