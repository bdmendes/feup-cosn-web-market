package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductQuantity struct {
	Product  primitive.ObjectID `json:"product_id"`
	Quantity int                `json:"quantity"`
}
