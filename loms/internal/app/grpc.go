package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/jackc/pgx/v5"
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

func NewGRPCApp(config config.Config, logger *slog.Logger) (*GRPCApp, error) {
	conn, err := dbConnect(context.Background(), config.Database.DSN)
	if err != nil {
		return nil, err
	}

	orderRepo := orderRepository.NewPostgresOrderRepository(conn)
	orderService := serviceOrder.NewOrderService(orderRepo)

	stocksRepo := repositoryStocks.NewPostgresStocksRepository(conn)
	stocksService := srviceStocks.NewStocksService(stocksRepo)

	orderHandler := handlerOrder.NewOrderHandler(orderService, stocksService)
	stocksHandler := handlerStocks.NewStocksHandler(stocksService)

	server := grpc.NewServer(getServerOption())

	loms.RegisterOrderServer(server, orderHandler)
	loms.RegisterStocksServer(server, stocksHandler)

	return &GRPCApp{
		cfg:    config,
		log:    logger,
		server: server,
	}, nil
}

func (a *GRPCApp) ListenAndServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port))
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

func dbConnect(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
