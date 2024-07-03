package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256/loms/internal/config"
	"route256/loms/internal/middleware"
	loms "route256/loms/pb/api"
)

type HTTPGW struct {
	cfg    config.Config
	log    *slog.Logger
	server *http.Server
	mux    *runtime.ServeMux
}

func headerMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-auth":
		return key, true
	default:
		return key, false
	}
}

func NewHTTPGW(cfg config.Config, log *slog.Logger) *HTTPGW {
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))

	return &HTTPGW{
		cfg:    cfg,
		log:    log,
		server: &http.Server{Addr: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port), Handler: middleware.WithHTTPLoggingMiddleware(mux, log)},
		mux:    mux,
	}
}

func (a *HTTPGW) ListenAndServe() error {
	conn, err := grpc.NewClient(fmt.Sprintf("dns:%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	a.mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	err = loms.RegisterOrderHandler(context.Background(), a.mux, conn)
	if err != nil {
		return err
	}

	err = loms.RegisterStocksHandler(context.Background(), a.mux, conn)
	if err != nil {
		return err
	}

	if err = a.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *HTTPGW) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
