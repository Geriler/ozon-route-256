package model

import (
	"route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

type Item struct {
	SKU   model.SKU
	Count int64
}

func LomsItemToItem(lomsItem *loms.Item) *Item {
	return &Item{
		SKU:   model.SKU(lomsItem.SkuId),
		Count: lomsItem.Count,
	}
}

func LomsItemsToItems(lomsItems []*loms.Item) []*Item {
	items := make([]*Item, 0)
	for _, item := range lomsItems {
		items = append(items, LomsItemToItem(item))
	}
	return items
}

func ItemToLomsItem(item *Item) *loms.Item {
	return &loms.Item{
		SkuId: int64(item.SKU),
		Count: item.Count,
	}
}

func ItemsToLomsItems(items []*Item) []*loms.Item {
	lomsItems := make([]*loms.Item, 0)
	for _, item := range items {
		lomsItems = append(lomsItems, ItemToLomsItem(item))
	}
	return lomsItems
}
