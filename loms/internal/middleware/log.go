package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"google.golang.org/grpc"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	const op = "middleware.Logger"

	logger := slog.With(
		slog.String("op", op),
		slog.String("method", info.FullMethod),
	)

	logger.Info("request received")

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return resp, err
}

func WithHTTPLoggingMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.WithHTTPLoggingMiddleware"

		logger := logger.With(
			slog.String("op", op),
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
		)

		logger.Info("request received")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
