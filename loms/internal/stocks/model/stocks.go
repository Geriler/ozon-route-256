package model

type SKU int64

type Stocks struct {
	SKU           SKU
	TotalCount    int64
	ReservedCount int64
}
