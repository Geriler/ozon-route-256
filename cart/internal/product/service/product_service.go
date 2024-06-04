package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	cartModel "route256/cart/internal/cart/model"
	"route256/cart/internal/product/model"
)

const baseUrl = "http://route256.pavl.uk:8080"

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) GetProduct(skuId cartModel.SkuID) (*model.Product, error) {
	url := baseUrl + "/get_product"

	request := model.GetProductRequest{
		Token: "testtoken",
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
