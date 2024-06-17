package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"sync"

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

var (
	ErrSkuNotFound    = errors.New("sku not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (r *InMemoryStocksRepository) Reserve(_ context.Context, sku model.SKU, count int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if stock, ok := r.stocks[sku]; ok {
		if stock.ReservedCount+count > stock.TotalCount {
			return ErrNotEnoughStock
		}

		stock.ReservedCount += count
		r.stocks[sku] = stock
		return nil
	}

	return ErrSkuNotFound
}

func (r *InMemoryStocksRepository) ReserveRemove(_ context.Context, sku model.SKU, count int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if stock, ok := r.stocks[sku]; ok {
		stock.ReservedCount -= count
		stock.TotalCount -= count
		r.stocks[sku] = stock
		return nil
	}

	return ErrSkuNotFound
}

func (r *InMemoryStocksRepository) ReserveCancel(_ context.Context, sku model.SKU, count int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if stock, ok := r.stocks[sku]; ok {
		stock.ReservedCount -= count
		r.stocks[sku] = stock
		return nil
	}

	return ErrSkuNotFound
}

func (r *InMemoryStocksRepository) GetBySKU(_ context.Context, sku model.SKU) (*model.Stocks, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if stock, ok := r.stocks[sku]; ok {
		return stock, nil
	}

	return nil, ErrSkuNotFound
}
