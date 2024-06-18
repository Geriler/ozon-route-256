package model

type CartResponse struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}

type CartCheckoutResponse struct {
	OrderID int64 `json:"order_id"`
}
