package handler

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/loms/internal"
	"route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

func (h *StocksHandler) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	ctx, span := h.tracer.Start(ctx, "StocksInfo", trace.WithAttributes(
		attribute.Int("sku_id", int(req.SkuId)),
	))
	defer span.End()

	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.StocksInfo", status)
	}(time.Now())

	span.AddEvent("Get stocks by SKU")
	stock, err := h.stocksService.GetBySKU(ctx, model.SKU(req.SkuId))
	if err != nil {
		status = "error"
		return nil, err
	}

	span.AddEvent("Return stocks info")
	return &loms.StocksInfoResponse{
		Count: stock.TotalCount - stock.ReservedCount,
	}, nil
}
