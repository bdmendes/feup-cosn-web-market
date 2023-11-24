package model

const (
	PaypalPaymentMethod     = "paypal"
	CreditCardPaymentMethod = "credit_card"
)

type PaymentRequestData struct {
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	PaymentData   string  `json:"payment_data" binding:"required"`
}
