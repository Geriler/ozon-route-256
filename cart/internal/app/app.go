package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"

	"route256/cart/internal/app/handler"
	cartHttp "route256/cart/internal/app/http"
	"route256/cart/internal/cart/repository"
	"route256/cart/internal/cart/service"
	"route256/cart/internal/config"
	lomsService "route256/cart/internal/loms/client"
	"route256/cart/internal/middleware"
	product "route256/cart/internal/product/service"
)

type App struct {
	mux     *http.ServeMux
	log     *slog.Logger
	server  *http.Server
	config  config.Config
	storage service.CartRepository
	loms    *lomsService.GRPCClient
	tracer  trace.Tracer
}

func NewApp(cfg config.Config, log *slog.Logger, loms *lomsService.GRPCClient, tracer trace.Tracer) *App {
	mux := http.NewServeMux()

	return &App{
		mux:     mux,
		log:     log,
		server:  &http.Server{Addr: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port), Handler: middleware.NewLogWrapperHandler(mux, log)},
		config:  cfg,
		storage: repository.NewInMemoryCartRepository(),
		loms:    loms,
		tracer:  tracer,
	}
}

func (a *App) ListenAndServe() error {
	productService := product.NewProductService(a.config.Product)
	cartService := service.NewCartService(a.storage)
	cartHandler := handler.NewCartHandler(cartService, productService, *a.loms, a.tracer)
	cartHttpHandlers := cartHttp.NewCartHttpHandlers(cartHandler, a.log)

	a.mux.Handle("GET /metrics", promhttp.Handler())
	a.mux.HandleFunc("POST /user/{user_id}/cart/checkout", cartHttpHandlers.Checkout)
	a.mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", cartHttpHandlers.AddItemsToCart)
	a.mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", cartHttpHandlers.DeleteItemsFromCart)
	a.mux.HandleFunc("DELETE /user/{user_id}/cart", cartHttpHandlers.DeleteCartByUserID)
	a.mux.HandleFunc("GET /user/{user_id}/cart/list", cartHttpHandlers.GetCartByUserID)

	if err := a.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
