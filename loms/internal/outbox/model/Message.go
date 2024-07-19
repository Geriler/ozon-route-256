package model

type Message struct {
	OrderID   int32  `json:"order_id"`
	EventType string `json:"event_type"`
}
