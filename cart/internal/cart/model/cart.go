package model

type UserID int64

type Cart struct {
	Items map[SkuID]*Item
}
