package routes

import (
	"cosn/consumers/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getAllNotificationsForConsumer(c *gin.Context) {
	notificatonsCollection := database.GetDatabase().Collection("notifications")

	id := c.Param("id")

	cursor, err := notificatonsCollection.Find(c, bson.M{"consumer": id})

	if err != nil {
		panic("Failed to get notifications: " + err.Error())
	}

	var notifications []bson.M
	err = cursor.All(c, &notifications)

	if err != nil {
		panic("Failed to get notifications: " + err.Error())
	}

	if len(notifications) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func registerForPriceNotification(c *gin.Context) {
	consumersCollection := database.GetDatabase().Collection("consumers")

	consumerId := c.Param("consumerId")

	productId := c.Param("productId")

	_, err := consumersCollection.UpdateOne(c, bson.M{"_id": consumerId},
		bson.M{"$addToSet": bson.M{"watchedProducts": productId}})

	if err != nil {
		panic("Failed to update consumer: " + err.Error())
	}

	c.JSON(http.StatusOK, gin.H{})
}

func AddNotificationRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET(":id", getAllNotificationsForConsumer)
	routerGroup.POST(":consumerId/products/:productId", registerForPriceNotification)
}
