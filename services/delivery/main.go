package main

import (
	"cosn/delivery/database"
	"cosn/delivery/observability"
	"cosn/delivery/routes"
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

	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
