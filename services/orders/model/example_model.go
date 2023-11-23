package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExampleModel struct {
	ID   primitive.ObjectID
	Name string
}

type ExampleModelRequestBody struct {
	Name string
}
