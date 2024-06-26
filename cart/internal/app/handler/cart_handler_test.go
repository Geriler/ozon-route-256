package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"route256/cart/internal/app/handler/mock"
	"route256/cart/internal/cart/model"
	"route256/cart/internal/cart/repository"
	"route256/cart/internal/loms/client"
	productModel "route256/cart/internal/product/model"
	loms "route256/cart/pb/api"
)

func TestCartHandler_AddItemsToCart(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	products := []*productModel.Product{
		{
			Name:  "Item 1",
			Price: 10,
		},
		{
			Name:  "Item 2",
			Price: 20,
		},
	}
	items := []model.Item{
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

	t.Run("add new item", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)
	})

	t.Run("add existing item", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		err = cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)
	})

	t.Run("add another item", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		productService.GetProductMock.Expect(items[1].SKU).Return(products[1], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[1].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[1].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[1])

		err = cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[1].SKU,
			Count:  items[1].Count,
		})
		require.Nil(t, err)
	})

	t.Run("add item with not found sku", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[2].SKU).Return(nil, errors.New("sku not found"))

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[2].SKU,
			Count:  items[2].Count,
		})
		require.EqualError(t, err, "sku not found")
	})

	t.Run("add item with not enough stocks", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: 0,
		}, nil)

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.ErrorIs(t, err, ErrNotEnoughStock)
	})
}

func TestCartHandler_DeleteItemsFromCart(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	products := []*productModel.Product{
		{
			Name:  "Item 1",
			Price: 10,
		},
		{
			Name:  "Item 2",
			Price: 20,
		},
	}
	items := []model.Item{
		{
			SKU:   1,
			Name:  products[0].Name,
			Count: 1,
			Price: products[0].Price,
		},
	}

	t.Run("delete not found item", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		cartService.DeleteItemsFromCartMock.Expect(context.Background(), userId, items[0].SKU)

		err := cartHandler.DeleteItemsFromCart(context.Background(), &model.UserSKURequest{
			UserID: userId,
			SKU:    items[0].SKU,
		})
		require.Nil(t, err)
	})

	t.Run("delete found item", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		cartService.DeleteItemsFromCartMock.Expect(context.Background(), userId, items[0].SKU)

		err = cartHandler.DeleteItemsFromCart(context.Background(), &model.UserSKURequest{
			UserID: userId,
			SKU:    items[0].SKU,
		})
		require.Nil(t, err)
	})
}

func TestCartHandler_DeleteCart(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	products := []*productModel.Product{
		{
			Name:  "Item 1",
			Price: 10,
		},
		{
			Name:  "Item 2",
			Price: 20,
		},
	}
	items := []model.Item{
		{
			SKU:   1,
			Name:  products[0].Name,
			Count: 1,
			Price: products[0].Price,
		},
	}

	t.Run("delete not found cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		cartService.DeleteCartByUserIDMock.Expect(context.Background(), userId)

		err := cartHandler.DeleteCart(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.Nil(t, err)
	})

	t.Run("delete empty cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		cartService.DeleteItemsFromCartMock.Expect(context.Background(), userId, items[0].SKU)

		err = cartHandler.DeleteItemsFromCart(context.Background(), &model.UserSKURequest{
			UserID: userId,
			SKU:    items[0].SKU,
		})
		require.Nil(t, err)

		cartService.DeleteCartByUserIDMock.Expect(context.Background(), userId)

		err = cartHandler.DeleteCart(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.Nil(t, err)
	})

	t.Run("delete not empty cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		cartService.DeleteCartByUserIDMock.Expect(context.Background(), userId)

		err = cartHandler.DeleteCart(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.Nil(t, err)
	})
}

func TestCartHandler_GetCart(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	products := []*productModel.Product{
		{
			Name:  "Item 1",
			Price: 10,
		},
		{
			Name:  "Item 2",
			Price: 20,
		},
	}
	items := []model.Item{
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
	}

	t.Run("get not found cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		cartService.GetCartByUserIDMock.Expect(context.Background(), userId).Return(&model.Cart{}, repository.ErrCartNotFoundOrEmpty)

		cartResponse, err := cartHandler.GetCart(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.ErrorIs(t, err, repository.ErrCartNotFoundOrEmpty)
		require.Equal(t, cartResponse, model.CartResponse{})
	})

	t.Run("get empty cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		cartService.DeleteItemsFromCartMock.Expect(context.Background(), userId, items[0].SKU)

		err = cartHandler.DeleteItemsFromCart(context.Background(), &model.UserSKURequest{
			UserID: userId,
			SKU:    items[0].SKU,
		})
		require.Nil(t, err)

		cartService.GetCartByUserIDMock.Expect(context.Background(), userId).Return(&model.Cart{}, repository.ErrCartNotFoundOrEmpty)

		cartResponse, err := cartHandler.GetCart(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.ErrorIs(t, err, repository.ErrCartNotFoundOrEmpty)
		require.Equal(t, cartResponse, model.CartResponse{})
	})

	t.Run("get not empty cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		cartItems := make(map[model.SkuID]*model.Item)
		cartItems[items[0].SKU] = &items[0]
		cart := &model.Cart{
			Items: cartItems,
		}

		productService.GetRPSMock.Expect().Return(10)
		cartService.GetCartByUserIDMock.Expect(context.Background(), userId).Return(cart, nil)
		cartService.GetTotalPriceMock.Expect(context.Background(), cart).Return(10)

		cartResponse, err := cartHandler.GetCart(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.Nil(t, err)
		require.EqualValues(t, uint32(10), cartResponse.TotalPrice)
		require.Len(t, cartResponse.Items, 1)
	})

	t.Run("get cart with difference items", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetRPSMock.Expect().Return(10)
		productService.GetProductMock.When(items[0].SKU).Then(products[0], nil).
			GetProductMock.When(items[1].SKU).Then(products[1], nil)

		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[1].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[1].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[1])

		err = cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[1].SKU,
			Count:  items[1].Count,
		})
		require.Nil(t, err)

		cartItems := make(map[model.SkuID]*model.Item)
		cartItems[items[0].SKU] = &items[0]
		cartItems[items[1].SKU] = &items[1]
		cart := &model.Cart{
			Items: cartItems,
		}

		cartService.GetCartByUserIDMock.Expect(context.Background(), userId).Return(cart, nil)
		cartService.GetTotalPriceMock.Expect(context.Background(), cart).Return(30)

		cartResponse, err := cartHandler.GetCart(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.Nil(t, err)
		require.EqualValues(t, uint32(30), cartResponse.TotalPrice)
		require.Len(t, cartResponse.Items, 2)
	})
}

func TestCartHandler_Checkout(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	products := []*productModel.Product{
		{
			Name:  "Item 1",
			Price: 10,
		},
	}
	items := []model.Item{
		{
			SKU:   1,
			Name:  products[0].Name,
			Count: 1,
			Price: products[0].Price,
		},
	}

	t.Run("checkout empty cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		cartService.GetCartByUserIDMock.Expect(context.Background(), userId).Return(&model.Cart{}, repository.ErrCartNotFoundOrEmpty)

		_, err := cartHandler.Checkout(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.ErrorIs(t, err, repository.ErrCartNotFoundOrEmpty)
	})

	t.Run("checkout not empty cart", func(t *testing.T) {
		t.Parallel()

		ctrl := minimock.NewController(t)
		productService := mock.NewProductServiceMock(ctrl)
		cartService := mock.NewCartServiceMock(ctrl)
		orderService := mock.NewOrderServiceMock(ctrl)
		stocksService := mock.NewStocksServiceMock(ctrl)
		grpcClient := client.NewGRPCClient(orderService, stocksService, nil)
		cartHandler := NewCartHandler(cartService, productService, *grpcClient)

		productService.GetProductMock.Expect(items[0].SKU).Return(products[0], nil)
		stocksService.StocksInfoMock.Expect(context.Background(), &loms.StocksInfoRequest{
			SkuId: int64(items[0].SKU),
		}).Return(&loms.StocksInfoResponse{
			Count: items[0].Count,
		}, nil)
		cartService.AddItemsToCartMock.Expect(context.Background(), userId, items[0])

		err := cartHandler.AddItemsToCart(context.Background(), &model.UserSKUCountRequest{
			UserID: userId,
			SKU:    items[0].SKU,
			Count:  items[0].Count,
		})
		require.Nil(t, err)

		cartItems := make(map[model.SkuID]*model.Item)
		cartItems[items[0].SKU] = &items[0]

		cartService.GetCartByUserIDMock.Expect(context.Background(), userId).Return(&model.Cart{
			Items: cartItems,
		}, nil)
		orderService.OrderCreateMock.Expect(context.Background(), &loms.OrderCreateRequest{
			UserId: int64(userId),
			Items: []*loms.Item{
				{
					SkuId: int64(items[0].SKU),
					Count: items[0].Count,
				},
			},
		}).Return(&loms.OrderCreateResponse{
			OrderId: 1,
		}, nil)
		cartService.DeleteCartByUserIDMock.Expect(context.Background(), userId)

		_, err = cartHandler.Checkout(context.Background(), &model.UserRequest{
			UserID: userId,
		})
		require.Nil(t, err)
	})
}
