package handler

import (
	"context"

	stocksModel "route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

type StocksService interface {
	StocksServiceReserve(ctx context.Context, sku stocksModel.SKU, count int64) error
	StocksServiceReserveRemove(ctx context.Context, sku stocksModel.SKU, count int64) error
	StocksServiceReserveCancel(ctx context.Context, sku stocksModel.SKU, count int64) error
	StocksServiceGetBySKU(ctx context.Context, sku stocksModel.SKU) (*stocksModel.Stocks, error)
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
