package main

import (
	"bytes"
	"context"
	"cosn/consumers/database"
	"cosn/consumers/model"
	"cosn/consumers/observability"
	"cosn/consumers/routes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func sendGetRequest(url string, callback func(*http.Response)) {
	fmt.Println("Sending GET request to", url)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

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

func populateProductsView() {
	sendGetRequest(os.Getenv("WAREHOUSE_PRODUCTS_SERVICE_URL"), func(resp *http.Response) {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var products []map[string]interface{}
		err = json.Unmarshal(b, &products)
		if err != nil {
			fmt.Println("Error unmarshalling response body:", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		productsCollection := database.GetDatabase().Collection("products")
		for _, product := range products {
			var productModel model.Product

			productModel.Name = product["name"].(string)
			productModel.ID = fmt.Sprintf("%f", product["id"].(float64))

			_, err = productsCollection.InsertOne(ctx, productModel)
			if err != nil {
				fmt.Println("Error inserting product:", err)
			} else {
				fmt.Println("Inserted product", productModel.Name)
			}
		}
	})
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	observability.AddHealthCheckRoutes(router)

	consumerRouterGroup := router.Group("/")
	routes.AddConsumersRoutes(consumerRouterGroup)

	notificationRouterGroup := router.Group("/notifications")
	routes.AddNotificationRoutes(notificationRouterGroup)

	recommendationRouterGroup := router.Group("/recommendations")
	routes.AddRecommendationRoutes(recommendationRouterGroup)

	shoppingBasketRouterGroup := router.Group("/basket")
	routes.AddShoppingBasketRoutes(shoppingBasketRouterGroup)

	return router
}

func main() {
	database.InitDatabase()
	defer database.DisconnectDatabase()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		fmt.Printf("Using default port %s\n", port)
	}

	populateProductsView()

	router := setupRouter()

	go routes.ProductsConsumer()

	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
