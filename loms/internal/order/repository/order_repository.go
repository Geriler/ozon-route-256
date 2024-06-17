package repository

import (
	"context"
	"fmt"
	"sync"

	"route256/loms/internal/order/model"
)

type InMemoryOrderRepository struct {
	orders map[model.OrderID]*model.Order
	mutex  *sync.RWMutex
}

var ErrOrderNotFound = fmt.Errorf("%s", "order not found")

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[model.OrderID]*model.Order),
		mutex:  &sync.RWMutex{},
	}
}

func (r *InMemoryOrderRepository) SetStatus(_ context.Context, orderID model.OrderID, status model.Status) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if order, ok := r.orders[orderID]; ok {
		order.Status = status
		return nil
	}

	return ErrOrderNotFound
}

func (r *InMemoryOrderRepository) GetOrder(_ context.Context, orderID model.OrderID) (*model.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if order, ok := r.orders[orderID]; ok {
		return order, nil
	}

	return nil, ErrOrderNotFound
}

func (r *InMemoryOrderRepository) Create(_ context.Context, order *model.Order) model.OrderID {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	orderID := model.OrderID(len(r.orders) + 1)
	r.orders[orderID] = order
	return orderID
}
