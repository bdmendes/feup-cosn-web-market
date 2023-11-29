package main

import (
	"cosn/delivery/database"
	"cosn/delivery/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	deliveryRouterGroup := router.Group("/")
	routes.AddDeliveryRoutes(deliveryRouterGroup)

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

	router.Run(":" + port)
}
