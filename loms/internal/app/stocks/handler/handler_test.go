package handler

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"route256/loms/internal/app/stocks/handler/mock"
	"route256/loms/internal/order/model"
	stocksModel "route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

func TestStocksHandler_StocksInfo(t *testing.T) {
	t.Parallel()

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

	t.Run("stocks info", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		stocksService := mock.NewStocksServiceMock(ctrl)
		stocksHandler := NewStocksHandler(stocksService, nil)

		stocksService.GetBySKUMock.Expect(context.Background(), items[0].SKU).Return(&stocksModel.Stocks{
			SKU:           items[0].SKU,
			TotalCount:    10,
			ReservedCount: 0,
		}, nil)

		_, err := stocksHandler.StocksInfo(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		})
		require.Nil(t, err)
	})

	t.Run("stocks info sku not found", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		stocksService := mock.NewStocksServiceMock(ctrl)
		stocksHandler := NewStocksHandler(stocksService, nil)

		stocksService.GetBySKUMock.Expect(context.Background(), items[0].SKU).Return(nil, stocksModel.ErrSkuNotFound)

		_, err := stocksHandler.StocksInfo(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		})
		require.ErrorIs(t, err, stocksModel.ErrSkuNotFound)
	})
}
