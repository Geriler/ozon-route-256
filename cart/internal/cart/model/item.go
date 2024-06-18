package model

type SkuID int64

type Item struct {
	SKU   SkuID
	Name  string
	Count int64
	Price uint32
}
