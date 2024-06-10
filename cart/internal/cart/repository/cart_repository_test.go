package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"route256/cart/internal/cart/model"
)

func TestInMemoryCartRepository_AddItems(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	items := []*model.Item{
		{
			SKU:   1,
			Name:  "item 1",
			Count: 1,
			Price: 1,
		},
		{
			SKU:   2,
			Name:  "item 2",
			Count: 1,
			Price: 1,
		},
	}

	t.Run("add new item", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.AddItems(userId, *items[0])
		cart = cartRepository.carts[userId]
		require.EqualValues(t, uint16(1), cart.Items[items[0].SKU].Count)
	})

	t.Run("add existing item", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.AddItems(userId, *items[0])
		cartRepository.AddItems(userId, *items[0])
		cart = cartRepository.carts[userId]
		require.EqualValues(t, uint16(2), cart.Items[items[0].SKU].Count)
	})

	t.Run("add another item", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.AddItems(userId, *items[0])
		cartRepository.AddItems(userId, *items[1])
		cart = cartRepository.carts[userId]
		require.EqualValues(t, uint16(1), cart.Items[items[0].SKU].Count)
		require.EqualValues(t, uint16(1), cart.Items[items[1].SKU].Count)
	})
}

func TestInMemoryCartRepository_DeleteItems(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	t.Run("delete not found item", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.DeleteItems(userId, item.SKU)
		cart = cartRepository.carts[userId]
		require.Nil(t, cart)
	})

	t.Run("delete found item", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.AddItems(userId, *item)

		cartRepository.DeleteItems(userId, item.SKU)
		cart = cartRepository.carts[userId]
		require.EqualValues(t, 0, len(cart.Items))
	})
}

func TestInMemoryCartRepository_DeleteCart(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	t.Run("delete not found cart", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.DeleteCart(userId)
		cart = cartRepository.carts[userId]
		require.Nil(t, cart)
	})

	t.Run("delete empty cart", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.AddItems(userId, *item)
		cartRepository.DeleteItems(userId, item.SKU)

		cartRepository.DeleteCart(userId)
		cart = cartRepository.carts[userId]
		require.Nil(t, cart)
	})

	t.Run("delete not empty cart", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var cart *model.Cart

		cartRepository.AddItems(userId, *item)

		cartRepository.DeleteCart(userId)
		cart = cartRepository.carts[userId]
		require.Nil(t, cart)
	})
}

func TestInMemoryCartRepository_GetCart(t *testing.T) {
	t.Parallel()

	var userId model.UserID = 1
	item := &model.Item{
		SKU:   1,
		Name:  "test",
		Count: 1,
		Price: 1,
	}

	t.Run("get not found cart", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var (
			cart *model.Cart
			err  error
		)

		cart, err = cartRepository.GetCart(userId)
		require.EqualError(t, err, ErrCartNotFoundOrEmpty)
		require.Nil(t, cart)
	})

	t.Run("get one item", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var (
			cart *model.Cart
			err  error
		)

		cartRepository.AddItems(userId, *item)

		cart, err = cartRepository.GetCart(userId)
		require.Nil(t, err)
		require.EqualValues(t, 1, len(cart.Items))
		require.EqualValues(t, 1, cart.Items[item.SKU].Count)
		require.EqualValues(t, 1, cart.Items[item.SKU].Price)
		require.EqualValues(t, "test", cart.Items[item.SKU].Name)
	})

	t.Run("get empty cart", func(t *testing.T) {
		t.Parallel()

		cartRepository := NewInMemoryCartRepository()

		var (
			cart *model.Cart
			err  error
		)

		cartRepository.AddItems(userId, *item)

		cartRepository.DeleteItems(userId, item.SKU)

		cart, err = cartRepository.GetCart(userId)
		require.EqualError(t, err, ErrCartNotFoundOrEmpty)
		require.Nil(t, cart)
	})
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
