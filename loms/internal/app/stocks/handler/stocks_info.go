package handler

import (
	"context"
	"time"

	"route256/loms/internal"
	"route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

func (h *StocksHandler) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.StocksInfo", status)
	}(time.Now())

	stock, err := h.stocksService.GetBySKU(ctx, model.SKU(req.SkuId))
	if err != nil {
		status = "error"
		return nil, err
	}

	return &loms.StocksInfoResponse{
		Count: stock.TotalCount - stock.ReservedCount,
	}, nil
}
