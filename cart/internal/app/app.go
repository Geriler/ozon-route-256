package app

import (
	"log/slog"
	"net/http"
	"strconv"

	"route256/cart/internal/app/handler"
	cartHttp "route256/cart/internal/app/http"
	"route256/cart/internal/cart/repository"
	"route256/cart/internal/cart/service"
	"route256/cart/internal/config"
	"route256/cart/internal/middleware"
	product "route256/cart/internal/product/service"
)

type App struct {
	mux     *http.ServeMux
	log     *slog.Logger
	server  *http.Server
	config  config.Config
	storage service.CartRepository
}

func NewApp(cfg config.Config, log *slog.Logger) *App {
	mux := http.NewServeMux()

	return &App{
		mux:     mux,
		log:     log,
		server:  &http.Server{Addr: "localhost:" + strconv.Itoa(cfg.Port), Handler: middleware.NewLogWrapperHandler(mux, log)},
		config:  cfg,
		storage: repository.NewInMemoryCartRepository(),
	}
}

func (a *App) ListenAndServe() error {
	productService := product.NewProductService(a.config.Product.BaseUrl, a.config.Product.Token)
	cartService := service.NewCartService(a.storage)
	cartHandler := handler.NewCartHandler(cartService, productService)
	cartHttpHandlers := cartHttp.NewCartHttpHandlers(cartHandler, a.log)

	a.mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", cartHttpHandlers.AddItemsToCart)
	a.mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", cartHttpHandlers.DeleteItemsFromCart)
	a.mux.HandleFunc("DELETE /user/{user_id}/cart", cartHttpHandlers.DeleteCartByUserID)
	a.mux.HandleFunc("GET /user/{user_id}/cart/list", cartHttpHandlers.GetCartByUserID)

	if err := a.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
