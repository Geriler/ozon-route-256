package handler

import (
	"context"

	"route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

func (h *StocksHandler) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	stock, err := h.stocksService.StocksServiceGetBySKU(ctx, model.SKU(req.SkuId))
	if err != nil {
		return nil, err
	}

	return &loms.StocksInfoResponse{
		Count: stock.TotalCount - stock.ReservedCount,
	}, nil
}
