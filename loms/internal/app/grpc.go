package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
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
	"route256/loms/pkg/infra/shards"
)

type GRPCApp struct {
	cfg    config.Config
	log    *slog.Logger
	server *grpc.Server
}

func NewGRPCApp(config config.Config, logger *slog.Logger) (*GRPCApp, error) {
	pools, err := dbConnect(context.Background(), config.Database.DSNs)
	if err != nil {
		return nil, err
	}

	sm := shards.NewManager(shards.GetMurmur3ShardFn(len(pools)), pools)

	conn, err := sm.GetShardByIndex(0)
	if err != nil {
		return nil, err
	}

	orderRepo := orderRepository.NewPostgresOrderRepository(sm, logger)
	orderService := serviceOrder.NewOrderService(orderRepo)

	stocksRepo := repositoryStocks.NewPostgresStocksRepository(conn, logger)
	stocksService := srviceStocks.NewStocksService(stocksRepo)

	orderHandler := handlerOrder.NewOrderHandler(orderService, stocksService)
	stocksHandler := handlerStocks.NewStocksHandler(stocksService)

	server := grpc.NewServer(getServerOption(),
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	)

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

func (a *GRPCApp) GracefulStop() {
	a.server.GracefulStop()
}

func getServerOption() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		middleware.GRPCSreWrapper,
		middleware.Logger,
	)
}

func dbConnect(ctx context.Context, dsns []string) ([]*pgxpool.Pool, error) {
	databases := make([]*pgxpool.Pool, 0, len(dsns))

	for _, dsn := range dsns {
		db, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, err
		}

		err = db.Ping(ctx)
		if err != nil {
			return nil, err
		}

		databases = append(databases, db)
	}

	return databases, nil
}
