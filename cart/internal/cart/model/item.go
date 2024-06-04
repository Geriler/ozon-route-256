package model

type SkuID int64

type Item struct {
	SKU   SkuID
	Name  string
	Count uint16
	Price uint32
}
