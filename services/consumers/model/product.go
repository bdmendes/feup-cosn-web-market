package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID
	Description string
	Prices      []float32
}

type ProductQuantity struct {
	Product  primitive.ObjectID
	Quantity int
}
