package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	cartModel "route256/cart/internal/cart/model"
	"route256/cart/internal/config"
	"route256/cart/internal/middleware"
	"route256/cart/internal/product/model"
	redisClient "route256/cart/pkg/infra/redis"
	"route256/cart/pkg/lib/tracing"
)

type ProductService struct {
	baseUrl  string
	token    string
	rpsLimit int
	redis    *redisClient.Client
	locks    map[string]*sync.Mutex
}

func NewProductService(cfg config.Config) *ProductService {
	client := redisClient.NewClient(cfg)

	return &ProductService{
		baseUrl:  cfg.Product.BaseUrl,
		token:    cfg.Product.Token,
		rpsLimit: cfg.Product.RPSLimit,
		redis:    client,
		locks:    make(map[string]*sync.Mutex),
	}
}

func (s *ProductService) GetProduct(ctx context.Context, skuId cartModel.SkuID) (*model.Product, error) {
	ctx, span := tracing.StartSpanFromContext(ctx, "GetProduct")
	defer span.End()

	statusCode := http.StatusOK

	defer func(createdAt time.Time) {
		middleware.ObserveRequestDurationSeconds(time.Since(createdAt).Seconds(), "POST /get_product", strconv.Itoa(statusCode))
	}(time.Now())

	cacheID := fmt.Sprintf("ProductService.GetProduct:%d", int64(skuId))

	_, exist := s.locks[cacheID]
	if !exist {
		s.locks[cacheID] = &sync.Mutex{}
	}
	s.locks[cacheID].Lock()
	defer s.locks[cacheID].Unlock()

	data, err := s.redis.Get(ctx, cacheID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if err == nil {
		var product model.Product
		err = json.Unmarshal([]byte(data), &product)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}

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

	err = s.redis.Set(ctx, cacheID, buf, s.redis.TTL)

	return &product, nil
}

func (s *ProductService) GetRPSLimit() int {
	return s.rpsLimit
}
