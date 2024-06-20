package service

import (
	"context"

	orderModel "route256/loms/internal/order/model"
	stocksModel "route256/loms/internal/stocks/model"
)

type StocksRepository interface {
	Reserve(ctx context.Context, items []*orderModel.Item) error
	ReserveRemove(ctx context.Context, items []*orderModel.Item) error
	ReserveCancel(ctx context.Context, items []*orderModel.Item) error
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

func (ss *StocksService) StocksServiceReserve(ctx context.Context, items []*orderModel.Item) error {
	return ss.stocksRepository.Reserve(ctx, items)
}

func (ss *StocksService) StocksServiceReserveRemove(ctx context.Context, items []*orderModel.Item) error {
	return ss.stocksRepository.ReserveRemove(ctx, items)
}

func (ss *StocksService) StocksServiceReserveCancel(ctx context.Context, items []*orderModel.Item) error {
	return ss.stocksRepository.ReserveCancel(ctx, items)
}

func (ss *StocksService) StocksServiceGetBySKU(ctx context.Context, sku stocksModel.SKU) (*stocksModel.Stocks, error) {
	return ss.stocksRepository.GetBySKU(ctx, sku)
}
