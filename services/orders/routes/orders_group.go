package routes

import (
	"cosn/orders/database"
	"cosn/orders/model"
	"net/http"

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
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	var orders []model.Order
	err = cursor.All(c, &orders)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
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

	ordersCollection := database.GetDatabase().Collection("orders")
	cursor, err := ordersCollection.Find(c, bson.M{"client": client_id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var orders []model.Order
	if err = cursor.All(c, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func createOrder(c *gin.Context) {
	var orderRequest model.NewOrderRequest
	if err := c.BindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	order := model.CreateOrderFromRequest(orderRequest)

	if err := model.IsOrderValid(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ordersCollection := database.GetDatabase().Collection("orders")

	if _, err := ordersCollection.InsertOne(c, order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
	//TODO: communicate with payment to authorize the purchase
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
			"interval_days": order.IntervalDays,
		},
	}

	doc := ordersCollection.FindOneAndUpdate(c, bson.M{"_id": order_id}, update)
	if doc.Err() != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
	//TODO: when authorized communicate with delivery service start delivery process
}

func AddOrdersRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/", getOrders)
	routerGroup.POST("/", createOrder)

	routerGroup.GET("/:id", getOrder)
	routerGroup.PUT("/:id", updateOrder)

	routerGroup.GET("/clients/:client_id", getClientOrders)
}
