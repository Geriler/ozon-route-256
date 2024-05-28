package model

type CartResponse struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
