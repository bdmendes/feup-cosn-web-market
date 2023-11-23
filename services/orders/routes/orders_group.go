package routes

import (
	"cosn/orders/database"
	"cosn/orders/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getOrder(c *gin.Context) {
	order_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	ordersCollection := database.GetDatabase().Collection("orders")

	var order model.Order
	err = ordersCollection.FindOne(c, bson.M{"_id": order_id}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func getOrders(c *gin.Context) {
	ordersCollection := database.GetDatabase().Collection("orders")

	cursor, err := ordersCollection.Find(c, bson.M{})

	if err != nil {
		panic("Failed to get example models: " + err.Error())
	}

	var orders []model.Order
	err = cursor.All(c, &orders)

	if err != nil {
		panic("Failed to get example models: " + err.Error())
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func getClientOrders(c *gin.Context) {
	client_id, err := primitive.ObjectIDFromHex(c.Param("client_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	println(client_id.Hex())

	ordersCollection := database.GetDatabase().Collection("orders")
	cursor, err := ordersCollection.Find(c, bson.M{"client": client_id})
	if err != nil {
		panic("Failed to get example models: " + err.Error())
	}

	var orders []model.Order
	err = cursor.All(c, &orders)
	if err != nil {
		panic("Failed to get example models: " + err.Error())
	}

	c.JSON(http.StatusOK, orders)
}

func createOrder(c *gin.Context) {
	var order model.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	status := model.PENDING
	order.Status = &status

	if err := model.IsOrderValid(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	order.Date = &(now)

	ordersCollection := database.GetDatabase().Collection("orders")

	if _, err := ordersCollection.InsertOne(c, order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func updateOrder(c *gin.Context) {
	order_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var order model.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ordersCollection := database.GetDatabase().Collection("orders")

	if err := model.IsOrderValid(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"status":        order.Status,
			"date":          order.Date,
			"interval_days": order.IntervalDays,
		},
	}

	doc := ordersCollection.FindOneAndUpdate(c, bson.M{"_id": order_id}, update)
	if doc.Err() != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func AddOrdersRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/", getOrders)
	routerGroup.POST("/", createOrder)

	routerGroup.GET("/:id", getOrder)
	routerGroup.PUT("/:id", updateOrder)

	routerGroup.GET("/clients/:client_id", getClientOrders)
}
