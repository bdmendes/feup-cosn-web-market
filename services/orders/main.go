package main

import (
	"cosn/orders/database"
	"cosn/orders/observability"
	"cosn/orders/routes"
	"cosn/orders/tasks"
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

	ordersRouterGroup := router.Group("/orders")
	routes.AddOrdersRoutes(ordersRouterGroup)

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

	go tasks.ScheduledOrdersTask()

	router := setupRouter()
	router.Run(":" + port)
}
