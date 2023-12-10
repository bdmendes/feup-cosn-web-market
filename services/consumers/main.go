package main

import (
	"cosn/consumers/database"
	"cosn/consumers/observability"
	"cosn/consumers/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

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

	router := setupRouter()

	go routes.ProductsConsumer()

	router.Run(":" + port)
}
