package routes

import (
	"context"
	"cosn/delivery/database"
	"cosn/delivery/model"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getDeliveryData(c *gin.Context) {
	deliveryCollection := database.GetDatabase().Collection("delivery")

	orderId, err := primitive.ObjectIDFromHex(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var delivery model.Delivery
	err = deliveryCollection.FindOne(context.Background(), bson.M{"order_id": orderId}).Decode(&delivery)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

func createDelivery(c *gin.Context) {
	deliveryCollection := database.GetDatabase().Collection("delivery")

	var deliveryRequestData model.DeliveryRequestData
	err := c.BindJSON(&deliveryRequestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	currentDateTime := time.Now()

	// Random delivery time between 10 minutes and 3 hours for all deliveries
	deliveryDateTime := currentDateTime.Add(time.Duration(float64(10+rand.Intn(170)) * time.Duration.Minutes(1)))

	if !deliveryRequestData.ExpressDelivery {
		// Add between 1 and 3 days for non-express deliveries
		deliveryDateTime = deliveryDateTime.Add(time.Duration(float64(1+rand.Intn(3)) * time.Duration.Hours(24)))
	}

	var delivery model.Delivery
	delivery.OrderID = deliveryRequestData.OrderID
	delivery.EstimatedDeliveryDateTime = deliveryDateTime.Format(time.RFC3339)
	delivery.Location = deliveryRequestData.Location

	_, err = deliveryCollection.InsertOne(context.Background(), delivery)

	if err != nil {
		panic(err)
	}

	c.Status(http.StatusOK)
}

func markDeliveryAsDone(c *gin.Context) {
	deliveryCollection := database.GetDatabase().Collection("delivery")

	orderId, err := primitive.ObjectIDFromHex(c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var delivery model.Delivery
	err = deliveryCollection.FindOne(context.Background(), bson.M{"order_id": orderId}).Decode(&delivery)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}

	if delivery.IsDone() {
		c.JSON(http.StatusNotModified, "Delivery already done")
		return
	}

	delivery.DeliveryDateTime = time.Now().Format(time.RFC3339)

	_, err = deliveryCollection.UpdateOne(context.Background(), bson.M{"order_id": orderId}, bson.M{"$set": delivery})

	if err != nil {
		panic(err)
	}

	c.Status(http.StatusOK)
}

func AddDeliveryRoutes(rg *gin.RouterGroup) {
	rg.GET("/:orderId", getDeliveryData)
	rg.POST("/", createDelivery)
	rg.POST("/:orderId/markAsDone", markDeliveryAsDone)
}
