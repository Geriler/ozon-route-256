package handler

import (
	"context"

	orderModel "route256/loms/internal/order/model"
	stocksModel "route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

type StocksService interface {
	Reserve(ctx context.Context, items []*orderModel.Item) error
	ReserveRemove(ctx context.Context, items []*orderModel.Item) error
	ReserveCancel(ctx context.Context, items []*orderModel.Item) error
	GetBySKU(ctx context.Context, sku stocksModel.SKU) (*stocksModel.Stocks, error)
}

type StocksHandler struct {
	loms.UnimplementedStocksServer
	stocksService StocksService
}

func NewStocksHandler(stocksService StocksService) *StocksHandler {
	return &StocksHandler{
		stocksService: stocksService,
	}
}
