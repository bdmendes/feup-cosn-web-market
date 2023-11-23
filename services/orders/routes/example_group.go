package routes

import (
	"cosn/orders/database"
	"cosn/orders/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getExampleModel(c *gin.Context) {
	exampleCollection := database.GetDatabase().Collection("example")

	cursor, err := exampleCollection.Find(c, bson.M{})

	if err != nil {
		panic("Failed to get example models: " + err.Error())
	}

	var exampleModels []model.ExampleModel
	err = cursor.All(c, &exampleModels)

	if err != nil {
		panic("Failed to get example models: " + err.Error())
	}

	if len(exampleModels) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, exampleModels)
}

func createExampleModel(c *gin.Context) {
	var requestBody model.ExampleModelRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	exampleCollection := database.GetDatabase().Collection("example")

	if _, err := exampleCollection.InsertOne(c, requestBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func AddExampleRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/model", getExampleModel)
	routerGroup.POST("/model", createExampleModel)
}
