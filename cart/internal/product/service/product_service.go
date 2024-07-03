package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"route256/cart/internal"
	cartModel "route256/cart/internal/cart/model"
	"route256/cart/internal/config"
	"route256/cart/internal/product/model"
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

func (s *ProductService) GetProduct(skuId cartModel.SkuID) (*model.Product, error) {
	statusCode := http.StatusOK

	defer func(createdAt time.Time) {
		internal.SaveMetrics(time.Since(createdAt).Seconds(), "POST /get_product", strconv.Itoa(statusCode))
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
