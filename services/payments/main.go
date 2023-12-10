package main

import (
	"cosn/payments/database"
	"cosn/payments/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	exampleRouterGroup := router.Group("/payment")
	routes.AddPaymentRoutes(exampleRouterGroup)

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
