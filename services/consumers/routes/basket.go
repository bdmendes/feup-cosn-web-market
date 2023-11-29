package routes

import (
	"cosn/consumers/database"
	"cosn/consumers/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getShoppingCart(c *gin.Context) {
	consumerId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumerCollection := database.GetDatabase().Collection("consumers")

	var consumer model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": consumerId}).Decode(&consumer); err != nil {
		panic("Failed to get consumer: " + err.Error())
	}

	if len(consumer.ShoppingCart) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, consumer.ShoppingCart)
}

func updateShoppingCart(c *gin.Context) {
	consumerId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumerCollection := database.GetDatabase().Collection("consumers")

	var consumer model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": consumerId}).Decode(&consumer); err != nil {
		panic("Failed to get consumer: " + err.Error())
	}

	var productQuantity model.ProductQuantity
	if err := c.ShouldBindJSON(&productQuantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumer.ShoppingCart = append(consumer.ShoppingCart, productQuantity)

	if _, err := consumerCollection.UpdateOne(c, bson.M{"_id": consumerId},
		bson.M{"$set": bson.M{"shoppingCart": consumer.ShoppingCart}}); err != nil {
		panic("Failed to update consumer: " + err.Error())
	}

	c.JSON(http.StatusOK, consumer.ShoppingCart)
}

func removeFromShoppingCart(c *gin.Context) {
	consumerId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumerCollection := database.GetDatabase().Collection("consumers")

	var consumer model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": consumerId}).Decode(&consumer); err != nil {
		panic("Failed to get consumer: " + err.Error())
	}

	var productQuantity model.ProductQuantity
	if err := c.ShouldBindJSON(&productQuantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newShoppingCart []model.ProductQuantity
	for _, product := range consumer.ShoppingCart {
		if product.Product != productQuantity.Product {
			newShoppingCart = append(newShoppingCart, product)
		}
	}

	if _, err := consumerCollection.UpdateOne(c, bson.M{"_id": consumerId},
		bson.M{"$set": bson.M{"shoppingCart": newShoppingCart}}); err != nil {
		panic("Failed to update consumer: " + err.Error())
	}

	c.JSON(http.StatusOK, newShoppingCart)
}

func AddShoppingBasketRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET(":id", getShoppingCart)
	routerGroup.PUT(":id", updateShoppingCart)
	routerGroup.DELETE(":id", removeFromShoppingCart)
}
