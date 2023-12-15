package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PriceDropNotification struct {
	ConsumerID primitive.ObjectID `bson:"consumerId" json:"consumerId"`
	ProductID  interface{}        `bson:"productId" json:"productId"`
	OldPrice   float32            `bson:"oldPrice,truncate" json:"oldPrice"`
	NewPrice   float32            `bson:"newPrice,truncate" json:"newPrice"`
}
