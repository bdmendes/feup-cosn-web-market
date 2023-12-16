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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		var updatedProductsView []model.ProductNotification
		err = json.Unmarshal(b, &updatedProductsView)
		if err != nil {
			fmt.Println("Error unmarshalling response body:", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		productsCollection := database.GetDatabase().Collection("products")

		for _, product := range updatedProductsView {
			var productModel model.Product
			if err = productsCollection.FindOne(ctx, bson.M{"id": product.ID}).Decode(&productModel); err == nil {
				if product.Name != "" {
					productModel.Name = product.Name
				}

				productModel.Category = product.Category
				productModel.Brand = product.Brand
				if productModel.Prices[len(productModel.Prices)-1] != product.Price {
					productModel.Prices = append(productModel.Prices, product.Price)
				}
			} else {
				productModel.ID = product.ID
				productModel.Name = product.Name
				productModel.Category = product.Category
				productModel.Brand = product.Brand
				productModel.Prices = []float32{product.Price}
			}

			opts := options.Update().SetUpsert(true)
			_, err = productsCollection.UpdateOne(ctx, bson.M{"id": product.ID}, bson.M{"$set": productModel}, opts)

			if err != nil {
				fmt.Println("Error inserting product:", err)
			} else {
				fmt.Println("Inserted product", product.Name)
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
