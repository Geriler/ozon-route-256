package handler

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"route256/loms/internal/app/order/handler/mock"
	"route256/loms/internal/order/model"
	modelStocks "route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

func TestOrderHandler_OrderCancel(t *testing.T) {
	t.Parallel()

	var (
		userID  int64 = 1
		orderID int64 = 1
	)
	items := []*model.Item{
		{
			SKU:   1,
			Count: 1,
		},
		{
			SKU:   2,
			Count: 1,
		},
	}

	t.Run("cancel order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderItems := make([]*model.Item, 1)
		orderItems[0] = items[0]

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(&model.Order{
			UserID: userID,
			Status: model.StatusAwaitingPayment,
			Items:  orderItems,
		}, nil)
		stocksService.ReserveCancelMock.Expect(context.Background(), orderItems).Return(nil)
		orderService.SetStatusMock.Expect(context.Background(), model.OrderID(orderID), model.StatusCanceled).Return(nil)

		_, err := orderHandler.OrderCancel(context.Background(), &loms.OrderCancelRequest{
			OrderId: orderID,
		})
		require.Nil(t, err)
	})

	t.Run("cancel not found order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderItems := make([]*model.Item, 1)
		orderItems[0] = items[0]

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(nil, model.ErrOrderNotFound)

		_, err := orderHandler.OrderCancel(context.Background(), &loms.OrderCancelRequest{
			OrderId: orderID,
		})
		require.ErrorIs(t, err, model.ErrOrderNotFound)
	})

	t.Run("cancel paid order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderItems := make([]*model.Item, 1)
		orderItems[0] = items[0]

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(&model.Order{
			UserID: userID,
			Status: model.StatusPaid,
			Items:  orderItems,
		}, nil)

		_, err := orderHandler.OrderCancel(context.Background(), &loms.OrderCancelRequest{
			OrderId: orderID,
		})
		require.ErrorIs(t, err, ErrOrderCannotCanceled)
	})
}

func TestOrderHandler_OrderCreate(t *testing.T) {
	t.Parallel()

	var (
		userID  int64 = 1
		orderID int64 = 1
	)
	items := []*model.Item{
		{
			OrderID: 0,
			SKU:     1,
			Count:   1,
		},
	}
	reserveItems := []*model.Item{
		{
			OrderID: 1,
			SKU:     1,
			Count:   1,
		},
	}
	lomsItems := []*loms.Item{
		{
			SkuId: int64(items[0].SKU),
			Count: items[0].Count,
		},
	}

	t.Run("create order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderService.CreateMock.Expect(context.Background(), &model.Order{
			UserID: userID,
			Status: model.StatusNew,
			Items:  items,
		}).Return(1, nil)
		stocksService.ReserveMock.Expect(context.Background(), reserveItems).Return(nil)
		orderService.SetStatusMock.Expect(context.Background(), model.OrderID(1), model.StatusAwaitingPayment).Return(nil)

		_, err := orderHandler.OrderCreate(context.Background(), &loms.OrderCreateRequest{
			UserId: userID,
			Items:  lomsItems,
		})
		require.Nil(t, err)
	})

	t.Run("create order not enough stocks", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderService.CreateMock.Expect(context.Background(), &model.Order{
			UserID: userID,
			Status: model.StatusNew,
			Items:  items,
		}).Return(1, nil)
		stocksService.ReserveMock.Expect(context.Background(), reserveItems).Return(modelStocks.ErrNotEnoughStock)
		orderService.SetStatusMock.Expect(context.Background(), model.OrderID(orderID), model.StatusFailed).Return(nil)

		_, err := orderHandler.OrderCreate(context.Background(), &loms.OrderCreateRequest{
			UserId: userID,
			Items:  lomsItems,
		})
		require.ErrorIs(t, err, modelStocks.ErrNotEnoughStock)
	})
}

func TestOrderHandler_OrderInfo(t *testing.T) {
	t.Parallel()

	var (
		orderID int64 = 1
	)
	items := []*model.Item{
		{
			SKU:   1,
			Count: 1,
		},
	}

	t.Run("get order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(&model.Order{
			UserID: 1,
			Status: model.StatusNew,
			Items:  items,
		}, nil)

		_, err := orderHandler.OrderInfo(context.Background(), &loms.OrderInfoRequest{
			OrderId: orderID,
		})
		require.Nil(t, err)
	})

	t.Run("not found order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(nil, model.ErrOrderNotFound)

		_, err := orderHandler.OrderInfo(context.Background(), &loms.OrderInfoRequest{
			OrderId: orderID,
		})
		require.ErrorIs(t, err, model.ErrOrderNotFound)
	})
}

func TestOrderHandler_OrderPay(t *testing.T) {
	t.Parallel()

	var (
		userID  int64 = 1
		orderID int64 = 1
	)
	items := []*model.Item{
		{
			SKU:   1,
			Count: 1,
		},
	}

	t.Run("pay order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(&model.Order{
			UserID: userID,
			Status: model.StatusAwaitingPayment,
			Items:  items,
		}, nil)
		stocksService.ReserveRemoveMock.Expect(context.Background(), items).Return(nil)
		orderService.SetStatusMock.Expect(context.Background(), model.OrderID(orderID), model.StatusPaid).Return(nil)

		_, err := orderHandler.OrderPay(context.Background(), &loms.OrderPayRequest{
			OrderId: orderID,
		})
		require.Nil(t, err)
	})

	t.Run("pay not found order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(nil, model.ErrOrderNotFound)

		_, err := orderHandler.OrderPay(context.Background(), &loms.OrderPayRequest{
			OrderId: orderID,
		})
		require.ErrorIs(t, err, model.ErrOrderNotFound)
	})

	t.Run("pay failed order", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		orderHandler := NewOrderHandler(orderService, stocksService, nil)

		orderService.GetOrderMock.Expect(context.Background(), model.OrderID(orderID)).Return(&model.Order{
			UserID: userID,
			Status: model.StatusFailed,
			Items:  items,
		}, nil)

		_, err := orderHandler.OrderPay(context.Background(), &loms.OrderPayRequest{
			OrderId: orderID,
		})
		require.ErrorIs(t, err, ErrOrderCannotPaid)
	})
}
