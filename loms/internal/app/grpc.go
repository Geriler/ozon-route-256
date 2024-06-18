package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	handlerOrder "route256/loms/internal/app/order/handler"
	handlerStocks "route256/loms/internal/app/stocks/handler"
	"route256/loms/internal/config"
	"route256/loms/internal/middleware"
	orderRepository "route256/loms/internal/order/repository"
	serviceOrder "route256/loms/internal/order/service"
	repositoryStocks "route256/loms/internal/stocks/repository"
	srviceStocks "route256/loms/internal/stocks/service"
	loms "route256/loms/pb/api"
)

type GRPCApp struct {
	cfg    config.Config
	log    *slog.Logger
	server *grpc.Server
}

func NewGRPCApp(config config.Config, logger *slog.Logger) *GRPCApp {
	server := grpc.NewServer(getServerOption())
	return &GRPCApp{
		cfg:    config,
		log:    logger,
		server: server,
	}
}

func (a *GRPCApp) ListenAndServe() error {
	orderRepo := orderRepository.NewInMemoryOrderRepository()
	orderService := serviceOrder.NewOrderService(orderRepo)

	stocksRepo, err := repositoryStocks.NewInMemoryStocksRepository()
	if err != nil {
		return err
	}
	stocksService := srviceStocks.NewStocksService(stocksRepo)

	orderHandler := handlerOrder.NewOrderHandler(orderService, stocksService)
	stocksHandler := handlerStocks.NewStocksHandler(stocksService)

	loms.RegisterOrderServer(a.server, orderHandler)
	loms.RegisterStocksServer(a.server, stocksHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
	if err != nil {
		return err
	}

	if err = a.server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func getServerOption() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		middleware.Logger,
	)
}
