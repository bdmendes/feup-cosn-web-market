package routes

import (
	"bytes"
	"context"
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

func SendPostRequest(payload []byte, url string, callback func(*http.Response, map[string]interface{}), data map[string]interface{}) {
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

	callback(resp, data)
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

func PaymentCallback(resp *http.Response, data map[string]interface{}) {
	if resp.StatusCode == http.StatusOK {
		order_id, err := primitive.ObjectIDFromHex(data["order_id"].(string))
		if err != nil {
			fmt.Println("Error: invalid order_id")
			return
		}

		ordersCollection := database.GetDatabase().Collection("orders")

		update := bson.M{"$set": bson.M{"status": model.AUTHORIZED}}

		doc := ordersCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": order_id}, update)
		if doc.Err() != nil {
			fmt.Println("Error: ", doc.Err())
			return
		}

		var order model.Order
		err = doc.Decode(&order)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		payload := []byte(`{"order_id": "` + order_id.Hex() +
			`", "location": "` + order.Location +
			`", "express_delivery": ` + strconv.FormatBool(order.ExpressDelivery) +
			`}`)

		go PublishOrder(order)
		go SendPostRequest(payload, os.Getenv("DELIVERY_SERVICE_URL"), deliveryCallback, data)
	} else {
		fmt.Println("Payment failed: ", resp.Status)
	}
}

func deliveryCallback(resp *http.Response, data map[string]interface{}) {
	if resp.StatusCode == http.StatusCreated {
		order_id, err := primitive.ObjectIDFromHex(data["order_id"].(string))
		if err != nil {
			fmt.Println("Error: invalid order_id")
			return
		}

		ordersCollection := database.GetDatabase().Collection("orders")

		update := bson.M{"$set": bson.M{"status": model.SHIPPED}}

		doc := ordersCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": order_id}, update)
		if doc.Err() != nil {
			fmt.Println("Error: ", doc.Err())
			return
		}
	} else {
		fmt.Println("Delivery failed: ", resp.Status)
	}
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

	doc, err := ordersCollection.InsertOne(c, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.Status(http.StatusCreated)

	order.ID = doc.InsertedID.(primitive.ObjectID)

	payload := []byte(`{"amount": ` + strconv.FormatFloat(order.PaymentData.Amount, 'f', -1, 64) +
		`, "payment_method": "` + order.PaymentData.PaymentMethod +
		`", "payment_data": "` + order.PaymentData.PaymentData +
		`"}`)
	go SendPostRequest(payload, os.Getenv("PAYMENT_SERVICE_URL")+"/payment", PaymentCallback, map[string]interface{}{"order_id": order.ID.Hex()})
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

	if updateOrderRequest.PaymentData != nil {
		updateSet["payment"] = *(updateOrderRequest.PaymentData)
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

	order := model.Order{}
	err = doc.Decode(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})

	if order.Status == model.PENDING && updateOrderRequest.PaymentData != nil {
		payload := []byte(`{"amount": ` + strconv.FormatFloat(order.PaymentData.Amount, 'f', -1, 64) +
			`, "payment_method": "` + order.PaymentData.PaymentMethod +
			`", "payment_data": "` + order.PaymentData.PaymentData +
			`"}`)
		go SendPostRequest(payload, os.Getenv("PAYMENT_SERVICE_URL")+"/payment", PaymentCallback, map[string]interface{}{"order_id": order_id.Hex()})
	} else if order.Status == model.AUTHORIZED && updateOrderRequest.Status != nil && *(updateOrderRequest.Status) == model.SHIPPED {
		payload := []byte(`{"order_id": "` + order_id.Hex() +
			`", "location": "` + order.Location +
			`", "express_delivery": ` + strconv.FormatBool(order.ExpressDelivery) +
			`}`)
		go SendPostRequest(payload, os.Getenv("DELIVERY_SERVICE_URL"), deliveryCallback, map[string]interface{}{"order_id": order_id.Hex()})
	}
}

func AddOrdersRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("", getOrders)
	routerGroup.POST("", createOrder)

	routerGroup.GET("/:id", getOrder)
	routerGroup.PUT("/:id", updateOrder)

	routerGroup.GET("/clients/:client_id", getClientOrders)
}
