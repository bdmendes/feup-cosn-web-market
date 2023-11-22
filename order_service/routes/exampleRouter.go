package routes

import (
	"orders-service/controllers"

	"github.com/gin-gonic/gin"
)

func addExampleRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", controllers.Ping)
}
