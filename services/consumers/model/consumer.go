package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Consumer struct {
	ID              primitive.ObjectID
	Name            string
	Location        string
	WatchedProducts []primitive.ObjectID
	ShoppingCart    []ProductQuantity
}

type ConsumerRequestBody struct {
	Name     string
	Location string
}
