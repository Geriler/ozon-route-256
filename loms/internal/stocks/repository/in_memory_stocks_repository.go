package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"sync"

	orderModel "route256/loms/internal/order/model"
	"route256/loms/internal/stocks/model"
)

//go:embed stocks.json
var stocksBytes []byte

type InMemoryStocksRepository struct {
	stocks map[model.SKU]*model.Stocks
	mutex  *sync.RWMutex
}

func NewInMemoryStocksRepository() (*InMemoryStocksRepository, error) {
	var stocksJson []*model.Stocks
	err := json.Unmarshal(stocksBytes, &stocksJson)
	if err != nil {
		return nil, err
	}

	stocks := make(map[model.SKU]*model.Stocks)

	for _, stock := range stocksJson {
		stocks[stock.SKU] = stock
	}

	return &InMemoryStocksRepository{
		stocks: stocks,
		mutex:  &sync.RWMutex{},
	}, nil
}

func (r *InMemoryStocksRepository) Reserve(_ context.Context, items []*orderModel.Item) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, item := range items {
		if stock, ok := r.stocks[item.SKU]; ok {
			if stock.ReservedCount+item.Count > stock.TotalCount {
				return model.ErrNotEnoughStock
			}
		} else {
			return model.ErrSkuNotFound
		}
	}

	for _, item := range items {
		r.stocks[item.SKU].ReservedCount += item.Count
	}

	return nil
}

func (r *InMemoryStocksRepository) ReserveRemove(_ context.Context, items []*orderModel.Item) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, item := range items {
		if _, ok := r.stocks[item.SKU]; !ok {
			return model.ErrSkuNotFound
		}
	}

	for _, item := range items {
		r.stocks[item.SKU].ReservedCount -= item.Count
		r.stocks[item.SKU].TotalCount -= item.Count
	}

	return nil
}

func (r *InMemoryStocksRepository) ReserveCancel(_ context.Context, items []*orderModel.Item) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, item := range items {
		if _, ok := r.stocks[item.SKU]; !ok {
			return model.ErrSkuNotFound
		}
	}

	for _, item := range items {
		r.stocks[item.SKU].ReservedCount -= item.Count
	}

	return nil
}

func (r *InMemoryStocksRepository) GetBySKU(_ context.Context, sku model.SKU) (*model.Stocks, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if stock, ok := r.stocks[sku]; ok {
		return stock, nil
	}

	return nil, model.ErrSkuNotFound
}
