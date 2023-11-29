package model

type PriceDropNotification struct {
	ConsumerID string  `bson:"consumerId" json:"consumerId"`
	ProductID  string  `bson:"productId" json:"productId"`
	OldPrice   float32 `bson:"oldPrice,truncate" json:"oldPrice"`
	NewPrice   float32 `bson:"newPrice,truncate" json:"newPrice"`
}
