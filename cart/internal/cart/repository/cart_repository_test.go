package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"route256/cart/internal/cart/model"
)

func TestInMemoryCartRepository_AddItems(t *testing.T) {
	cartRepository := NewInMemoryCartRepository()

	var (
		userId model.UserID = 1
		cart   *model.Cart
	)
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	cart = cartRepository.carts[userId]
	require.Nil(t, cart)

	cartRepository.AddItems(userId, *item)
	cart = cartRepository.carts[userId]
	require.EqualValues(t, 1, cart.Items[item.SKU].Count)

	cartRepository.AddItems(userId, *item)
	cartRepository.AddItems(userId, *item)
	cart = cartRepository.carts[userId]
	require.EqualValues(t, 3, cart.Items[item.SKU].Count)
}

func TestInMemoryCartRepository_DeleteItems(t *testing.T) {
	cartRepository := NewInMemoryCartRepository()

	var (
		userId model.UserID = 1
		cart   *model.Cart
	)
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	cartRepository.DeleteItems(userId, item.SKU)
	cart = cartRepository.carts[userId]
	require.Nil(t, cart)

	cartRepository.AddItems(userId, *item)

	cartRepository.DeleteItems(userId, item.SKU)
	cart = cartRepository.carts[userId]
	require.EqualValues(t, 0, len(cart.Items))
}

func TestInMemoryCartRepository_DeleteCart(t *testing.T) {
	cartRepository := NewInMemoryCartRepository()

	var (
		userId model.UserID = 1
		cart   *model.Cart
	)
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	cartRepository.DeleteCart(userId)
	cart = cartRepository.carts[userId]
	require.Nil(t, cart)

	cartRepository.AddItems(userId, *item)

	cartRepository.DeleteCart(userId)
	cart = cartRepository.carts[userId]
	require.Nil(t, cart)
}

func TestInMemoryCartRepository_GetCart(t *testing.T) {
	cartRepository := NewInMemoryCartRepository()

	var (
		userId model.UserID = 1
		cart   *model.Cart
		err    error
	)
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	cart, err = cartRepository.GetCart(userId)
	require.EqualError(t, err, ErrCartNotFoundOrEmpty)
	require.Nil(t, cart)

	cartRepository.AddItems(userId, *item)

	cart, err = cartRepository.GetCart(userId)
	require.Nil(t, err)
	require.EqualValues(t, 1, len(cart.Items))

	cartRepository.DeleteItems(userId, item.SKU)

	cart, err = cartRepository.GetCart(userId)
	require.EqualError(t, err, ErrCartNotFoundOrEmpty)
	require.Nil(t, cart)
}

func BenchmarkInMemoryCartRepository_AddItems(b *testing.B) {
	cartRepository := NewInMemoryCartRepository()

	var userId model.UserID = 1
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	for i := 0; i < b.N; i++ {
		cartRepository.AddItems(userId, *item)
	}
}
