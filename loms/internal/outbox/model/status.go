package model

type Status string

const (
	StatusPending Status = "pending"
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)
