package model

import (
	"cosn/consumers/database"
	"sort"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	var latestProducts []Product
	opts := options.Find().SetLimit(20).SetSort(bson.M{"_id": -1})
	cursor, err := database.GetDatabase().Collection("products").Find(c, bson.M{}, opts)
	if err != nil {
		panic("Failed to get products: " + err.Error())
	}
	if err := cursor.All(c, &latestProducts); err != nil {
		panic("Failed to get products: " + err.Error())
	}

	sort.Slice(latestProducts, func(i, j int) bool {
		return latestProducts[i].SimilarityMultiple(products) > latestProducts[j].SimilarityMultiple(products)
	})

	products = append(products, latestProducts...)

	return products
}
