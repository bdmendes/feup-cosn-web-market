package routes

import (
	"cosn/consumers/database"
	"cosn/consumers/model"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAllNotificationsForConsumer(c *gin.Context) {
	notificatonsCollection := database.GetDatabase().Collection("notifications")

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursor, err := notificatonsCollection.Find(c, bson.M{"consumerId": id})

	if err != nil {
		panic("Failed to get notifications: " + err.Error())
	}

	var notifications []model.PriceDropNotification
	err = cursor.All(c, &notifications)

	if err != nil {
		panic("Failed to get notifications: " + err.Error())
	}

	if len(notifications) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func registerForPriceNotification(c *gin.Context) {
	consumersCollection := database.GetDatabase().Collection("consumers")

	consumerId, err := primitive.ObjectIDFromHex(c.Param("consumerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productId := c.Param("productId")

	_, err = consumersCollection.UpdateOne(c, bson.M{"_id": consumerId},
		bson.M{"$addToSet": bson.M{"watchedProducts": productId}})

	if err != nil {
		panic("Failed to update consumer: " + err.Error())
	}

	c.Status(http.StatusOK)
}

func AddNotificationRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET(":id", getAllNotificationsForConsumer)
	routerGroup.POST(":consumerId/products/:productId", registerForPriceNotification)
}
