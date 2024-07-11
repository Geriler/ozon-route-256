package handler

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

func (h *StocksHandler) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	ctx, span := h.tracer.Start(ctx, "StocksInfo", trace.WithAttributes(
		attribute.Int("sku_id", int(req.SkuId)),
	))
	defer span.End()

	span.AddEvent("Get stocks by SKU")
	stock, err := h.stocksService.GetBySKU(ctx, model.SKU(req.SkuId))
	if err != nil {
		return nil, err
	}

	span.AddEvent("Return stocks info")
	return &loms.StocksInfoResponse{
		Count: stock.TotalCount - stock.ReservedCount,
	}, nil
}
