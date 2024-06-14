package model

type Status string

const (
	StatusNew             Status = "new"
	StatusAwaitingPayment Status = "awaiting_payment"
	StatusPayed           Status = "payed"
	StatusFailed          Status = "failed"
	StatusCanceled        Status = "canceled"
)

type OrderID int64

type Order struct {
	UserID int64
	Status Status
	Items  []*Item
}
