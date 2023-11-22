package routes

import (
	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	exampleRoutes := router.Group("/example")
	addExampleRoutes(exampleRoutes)

	return router
}
