package handler

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"route256/cart/internal/app/handler/mock"
	"route256/cart/internal/cart/model"
	"route256/cart/internal/cart/repository"
	productModel "route256/cart/internal/product/model"
)

func TestCartHandler(t *testing.T) {
	ctrl := minimock.NewController(t)
	productService := mock.NewProductServiceMock(ctrl)
	cartService := mock.NewCartServiceMock(ctrl)

	var (
		err          error
		userId       = model.UserID(1)
		cartResponse = model.CartResponse{}
		cartItems    []model.Item
		cart         = &model.Cart{}
		products     = []*productModel.Product{
			{
				Name:  "Item 1",
				Price: 10,
			},
			{
				Name:  "Item 2",
				Price: 20,
			},
		}
		items = []model.Item{
			{
				SKU:   1,
				Name:  products[0].Name,
				Count: 1,
				Price: products[0].Price,
			},
			{
				SKU:   2,
				Name:  products[1].Name,
				Count: 1,
				Price: products[1].Price,
			},
			{
				SKU:   3,
				Count: 1,
			},
		}
	)

	cartHandler := NewCartHandler(cartService, productService)

	cartService.GetCartByUserIDMock.Expect(userId).Return(nil, errors.New(repository.ErrCartNotFoundOrEmpty))

	cartResponse, err = cartHandler.GetCart(&model.UserRequest{
		UserID: userId,
	})
	require.EqualError(t, err, repository.ErrCartNotFoundOrEmpty)
	require.Equal(t, cartResponse, model.CartResponse{})

	productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
	cartService.AddItemsToCartMock.Expect(userId, items[0])

	cart = &model.Cart{
		Items: map[model.SkuID]*model.Item{
			items[0].SKU: &items[0],
		},
	}
	cartItems = append(cartItems, items[0])
	cartService.GetCartByUserIDMock.Expect(userId).Return(cart, nil)
	cartService.GetTotalPriceMock.Expect(cart).Return(uint32(10))

	err = cartHandler.AddItemsToCart(&model.UserSKUCountRequest{
		UserID: userId,
		SKU:    items[0].SKU,
		Count:  items[0].Count,
	})

	require.NoError(t, err)

	cartResponse, err = cartHandler.GetCart(&model.UserRequest{
		UserID: userId,
	})
	require.NoError(t, err)
	require.Equal(t, cartResponse.Items, cartItems)
	require.Equal(t, cartResponse.TotalPrice, uint32(10))

	err = cartHandler.AddItemsToCart(&model.UserSKUCountRequest{
		UserID: userId,
		SKU:    items[0].SKU,
		Count:  items[0].Count,
	})

	require.NoError(t, err)

	cartItems[0].Count = 2
	cart.Items[items[0].SKU].Count = 2
	cartService.GetCartByUserIDMock.Expect(userId).Return(cart, nil)
	cartService.GetTotalPriceMock.Expect(cart).Return(uint32(20))

	cartResponse, err = cartHandler.GetCart(&model.UserRequest{
		UserID: userId,
	})
	require.NoError(t, err)
	require.Equal(t, cartResponse.Items, cartItems)
	require.Equal(t, cartResponse.TotalPrice, uint32(20))

	productService.GetProductMock.Expect(items[1].SKU).Return(products[1], nil)

	cartItems = append(cartItems, items[1])
	cart.Items[items[1].SKU] = &items[1]
	cartService.GetCartByUserIDMock.Expect(userId).Return(cart, nil)
	cartService.GetTotalPriceMock.Expect(cart).Return(uint32(40))

	cartService.AddItemsToCartMock.Expect(userId, items[1])

	err = cartHandler.AddItemsToCart(&model.UserSKUCountRequest{
		UserID: userId,
		SKU:    items[1].SKU,
		Count:  items[1].Count,
	})

	require.NoError(t, err)

	cartResponse, err = cartHandler.GetCart(&model.UserRequest{
		UserID: userId,
	})
	require.NoError(t, err)
	require.Equal(t, cartResponse.Items, cartItems)
	require.Equal(t, cartResponse.TotalPrice, uint32(40))

	productService.GetProductMock.Expect(items[2].SKU).Return(nil, errors.New("sku not found"))

	err = cartHandler.AddItemsToCart(&model.UserSKUCountRequest{
		UserID: userId,
		SKU:    items[2].SKU,
		Count:  items[2].Count,
	})
	require.EqualError(t, err, "sku not found")

	cartResponse, err = cartHandler.GetCart(&model.UserRequest{
		UserID: userId,
	})
	require.NoError(t, err)
	require.Equal(t, cartResponse.Items, cartItems)
	require.Equal(t, cartResponse.TotalPrice, uint32(40))

	cartService.DeleteItemsFromCartMock.Expect(userId, items[0].SKU)
	cartItems = append(cartItems[1:])
	delete(cart.Items, items[0].SKU)
	cartService.GetCartByUserIDMock.Expect(userId).Return(cart, nil)
	cartService.GetTotalPriceMock.Expect(cart).Return(uint32(20))

	err = cartHandler.DeleteItemsFromCart(&model.UserSKURequest{
		UserID: userId,
		SKU:    items[0].SKU,
	})
	require.NoError(t, err)

	cartResponse, err = cartHandler.GetCart(&model.UserRequest{
		UserID: userId,
	})
	require.NoError(t, err)
	require.Equal(t, cartResponse.Items, cartItems)
	require.Equal(t, cartResponse.TotalPrice, uint32(20))

	cartService.DeleteCartByUserIDMock.Expect(userId)
	cartService.GetCartByUserIDMock.Expect(userId).Return(nil, errors.New(repository.ErrCartNotFoundOrEmpty))

	err = cartHandler.DeleteCart(&model.UserRequest{
		UserID: userId,
	})
	require.NoError(t, err)

	cartResponse, err = cartHandler.GetCart(&model.UserRequest{
		UserID: userId,
	})
	require.EqualError(t, err, repository.ErrCartNotFoundOrEmpty)
}
