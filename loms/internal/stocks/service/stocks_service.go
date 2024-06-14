package service

import (
	"context"

	stocksModel "route256/loms/internal/stocks/model"
)

type StocksRepository interface {
	Reserve(ctx context.Context, sku stocksModel.SKU, count int64) error
	ReserveRemove(ctx context.Context, sku stocksModel.SKU, count int64) error
	ReserveCancel(ctx context.Context, sku stocksModel.SKU, count int64) error
	GetBySKU(ctx context.Context, sku stocksModel.SKU) (*stocksModel.Stocks, error)
}

type StocksService struct {
	stocksRepository StocksRepository
}

func NewStocksService(stocksRepository StocksRepository) *StocksService {
	return &StocksService{
		stocksRepository: stocksRepository,
	}
}

func (ss *StocksService) StocksServiceReserve(ctx context.Context, sku stocksModel.SKU, count int64) error {
	return ss.stocksRepository.Reserve(ctx, sku, count)
}

func (ss *StocksService) StocksServiceReserveRemove(ctx context.Context, sku stocksModel.SKU, count int64) error {
	return ss.stocksRepository.ReserveRemove(ctx, sku, count)
}

func (ss *StocksService) StocksServiceReserveCancel(ctx context.Context, sku stocksModel.SKU, count int64) error {
	return ss.stocksRepository.ReserveCancel(ctx, sku, count)
}

func (ss *StocksService) StocksServiceGetBySKU(ctx context.Context, sku stocksModel.SKU) (*stocksModel.Stocks, error) {
	return ss.stocksRepository.GetBySKU(ctx, sku)
}
