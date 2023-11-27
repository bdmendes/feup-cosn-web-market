package routes

import (
	"bytes"
	"cosn/orders/database"
	"cosn/orders/model"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func sendPostRequest(payload []byte, url string, callback func(*http.Response)) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	callback(resp)
}

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

func paymentCallback(resp *http.Response) {
	fmt.Println("Payment callback")
	//TODO: update order status to authorized
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

	payload := []byte(`{"amount": ` + strconv.FormatFloat(order.TotalPrice, 'f', -1, 64) + `, "payment_method": "paypal", "payment_data": "` + order.Payment + `"}`)
	go sendPostRequest(payload, os.Getenv("PAYMENT_SERVICE_URL")+"/payment", paymentCallback)
}

func deliveryCallback(resp *http.Response) {
	fmt.Println("Delivery callback")
	//TODO: update order status to authorized
}

func updateOrder(c *gin.Context) {
	order_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var updateOrderRequest model.UpdateOrderRequest
	if err := c.BindJSON(&updateOrderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ordersCollection := database.GetDatabase().Collection("orders")

	updateSet := bson.M{}

	if updateOrderRequest.Location != nil {
		updateSet["location"] = *(updateOrderRequest.Location)
	}

	if updateOrderRequest.Payment != nil {
		updateSet["payment"] = *(updateOrderRequest.Payment)
	}

	if updateOrderRequest.Status != nil {
		updateSet["status"] = *(updateOrderRequest.Status)
	}

	if updateOrderRequest.IntervalDays != nil {
		updateSet["intervaldays"] = *(updateOrderRequest.IntervalDays)
	}

	update := bson.M{"$set": updateSet}

	doc := ordersCollection.FindOneAndUpdate(c, bson.M{"_id": order_id}, update)
	if doc.Err() != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
	//TODO: when authorized communicate with delivery service start delivery process

	payload := []byte(`{}`)
	go sendPostRequest(payload, os.Getenv("DELIVERY_SERVICE_URL")+"/delivery", deliveryCallback)
}

func AddOrdersRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/", getOrders)
	routerGroup.POST("/", createOrder)

	routerGroup.GET("/:id", getOrder)
	routerGroup.PUT("/:id", updateOrder)

	routerGroup.GET("/clients/:client_id", getClientOrders)
}
