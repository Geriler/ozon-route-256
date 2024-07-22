package model

import orderModel "route256/loms/internal/order/model"

type Message struct {
	OrderID   int32             `json:"order_id"`
	EventType orderModel.Status `json:"event_type"`
}
