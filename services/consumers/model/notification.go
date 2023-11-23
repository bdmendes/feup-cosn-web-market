package model

type PriceDropNotification struct {
	ConsumerID string
	ProductID  string
	OldPrice   float32
	NewPrice   float32
}
