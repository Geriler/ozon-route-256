package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	cartModel "route256/cart/internal/cart/model"
	"route256/cart/internal/config"
	"route256/cart/internal/middleware"
	"route256/cart/internal/product/model"
	"route256/cart/pkg/lib/tracing"
)

type ProductService struct {
	baseUrl  string
	token    string
	rpsLimit int
}

func NewProductService(cfg config.ProductConfig) *ProductService {
	return &ProductService{
		baseUrl:  cfg.BaseUrl,
		token:    cfg.Token,
		rpsLimit: cfg.RPSLimit,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, skuId cartModel.SkuID) (*model.Product, error) {
	ctx, span := tracing.StartSpanFromContext(ctx, "GetProduct")
	defer span.End()

	statusCode := http.StatusOK

	defer func(createdAt time.Time) {
		middleware.ObserveRequestDurationSeconds(time.Since(createdAt).Seconds(), "POST /get_product", strconv.Itoa(statusCode))
	}(time.Now())

	url := s.baseUrl + "/get_product"

	request := model.GetProductRequest{
		Token: s.token,
		Sku:   skuId,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	statusCode = response.StatusCode

	if response.StatusCode != http.StatusOK {
		var errorResponse model.GetProductErrorResponse
		err = json.Unmarshal(buf, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.Message)
	}

	var product model.Product
	err = json.Unmarshal(buf, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *ProductService) GetRPSLimit() int {
	return s.rpsLimit
}
