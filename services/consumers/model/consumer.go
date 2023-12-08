package model

import (
	"cosn/consumers/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Consumer struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Name            string
	Location        string
	WatchedProducts []primitive.ObjectID
	ShoppingCart    []ProductQuantity
	OrderHistory    []ProductQuantity
}

type ConsumerRequestBody struct {
	Name     string
	Location string
}

func (consumer *Consumer) RelatedProducts(c *gin.Context) []Product {
	var products []Product

	for _, productQuantity := range consumer.ShoppingCart {
		var product Product
		if err := database.GetDatabase().Collection("products").FindOne(c,
			bson.M{"_id": productQuantity.Product}).Decode(&product); err != nil {
			panic("Failed to get product: " + err.Error())
		}
		products = append(products, product)
	}

	for _, productQuantity := range consumer.OrderHistory {
		var product Product
		if err := database.GetDatabase().Collection("products").FindOne(c,
			bson.M{"_id": productQuantity.Product}).Decode(&product); err != nil {
			panic("Failed to get product: " + err.Error())
		}
		products = append(products, product)
	}

	for _, watchedProduct := range consumer.WatchedProducts {
		var product Product
		if err := database.GetDatabase().Collection("products").FindOne(c,
			bson.M{"_id": watchedProduct}).Decode(&product); err != nil {
			panic("Failed to get product: " + err.Error())
		}

		products = append(products, product)
	}

	return products
}
