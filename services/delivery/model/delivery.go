package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Delivery struct {
	ID                        primitive.ObjectID
	OrderID                   int64  `json:"order_id"`
	EstimatedDeliveryDateTime string `json:"estimated_delivery_datetime"`
	DeliveryDateTime          string `json:"delivery_datetime"`
}

func (d *Delivery) IsDone() bool {
	return d.DeliveryDateTime != ""
}
