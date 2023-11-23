package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PENDING    = "PENDING"
	AUTHORIZED = "AUTHORIZED"
	DELIVERED  = "DELIVERED"
	CANCELLED  = "CANCELLED"
)

type Order struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Client       primitive.ObjectID `json:"client_id"`
	Description  *string            `json:"description"`
	Location     string             `json:"location"`
	Products     []ProductQuantity  `json:"products"`
	TotalPrice   float64            `json:"total_price"`
	Date         *time.Time         `json:"date"`
	IntervalDays *int64             `json:"interval_days"`
	Payment      string             `json:"payment"`
	Status       *string            `json:"status"`
}

func IsOrderValid(order Order) error {
	if order.IntervalDays != nil && *(order.IntervalDays) < 1 {
		return errors.New("invalid interval")
	}

	if order.TotalPrice < 0 {
		return errors.New("invalid price")
	}

	if len(order.Products) <= 0 {
		return errors.New("order with no products")
	}

	if *(order.Status) != PENDING && *(order.Status) != AUTHORIZED &&
		*(order.Status) != DELIVERED && *(order.Status) != CANCELLED {
		return errors.New("invalid order status")
	}

	for i := 0; i < len(order.Products); i++ {
		if order.Products[i].Quantity < 1 {
			return errors.New("product with invalid quantity")
		}
	}

	return nil
}
