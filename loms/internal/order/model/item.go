package model

import (
	"route256/loms/internal/stocks/model"
	loms "route256/loms/pb/api"
)

type Item struct {
	OrderID int64
	SKU     model.SKU
	Count   int64
}

func LomsItemToItem(lomsItem *loms.Item) *Item {
	if lomsItem == nil {
		return nil
	}
	return &Item{
		SKU:   model.SKU(lomsItem.SkuId),
		Count: lomsItem.Count,
	}
}

func LomsItemsToItems(lomsItems []*loms.Item) []*Item {
	items := make([]*Item, len(lomsItems))
	for i, item := range lomsItems {
		items[i] = LomsItemToItem(item)
	}
	return items
}

func ItemToLomsItem(item *Item) *loms.Item {
	if item == nil {
		return nil
	}
	return &loms.Item{
		SkuId: int64(item.SKU),
		Count: item.Count,
	}
}

func ItemsToLomsItems(items []*Item) []*loms.Item {
	lomsItems := make([]*loms.Item, len(items))
	for i, item := range items {
		lomsItems[i] = ItemToLomsItem(item)
	}
	return lomsItems
}
