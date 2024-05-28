package middleware

import (
	"log/slog"
	"net/http"
)

type LogWrapperHandler struct {
	wrap   http.Handler
	logger *slog.Logger
}

func NewLogWrapperHandler(wrap http.Handler, logger *slog.Logger) *LogWrapperHandler {
	return &LogWrapperHandler{
		wrap:   wrap,
		logger: logger,
	}
}

func (h LogWrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "middleware.LogWrapperHandler.ServeHTTP"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("url", r.URL.String()),
		slog.String("method", r.Method),
	)

	log.Info("request received")

	h.wrap.ServeHTTP(w, r)
}
