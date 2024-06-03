package app

import (
	"log/slog"
	"net"
	"net/http"

	"route256/cart/internal/middleware"
)

type App struct {
	mux *http.ServeMux
	log *slog.Logger
}

func NewApp(mux *http.ServeMux, log *slog.Logger) *App {
	return &App{
		mux: mux,
		log: log,
	}
}

func (a *App) ListenAndServe() error {
	conn, err := net.Listen("tcp", ":8082")
	if err != nil {
		return err
	}
	defer conn.Close()

	logWrapperHandler := middleware.NewLogWrapperHandler(a.mux, a.log)

	if err = http.Serve(conn, logWrapperHandler); err != nil {
		return err
	}
	return nil
}
