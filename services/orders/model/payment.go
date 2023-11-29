package model

type PaymentData struct {
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	PaymentData   string  `json:"payment_data" binding:"required"`
}
