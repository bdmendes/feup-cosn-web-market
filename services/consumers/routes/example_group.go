package routes

import (
	"cosn/template/database"
	"cosn/template/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func createConsumer(c *gin.Context) {
	var consumerRequestBody model.ConsumerRequestBody
	if err := c.ShouldBindJSON(&consumerRequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumerCollection := database.GetDatabase().Collection("consumers")

	consumerModel := model.Consumer{
		Name:     consumerRequestBody.Name,
		Location: consumerRequestBody.Location,
	}

	if _, err := consumerCollection.InsertOne(c, consumerModel); err != nil {
		panic("Failed to create consumer model: " + err.Error())
	}

	c.JSON(http.StatusCreated, consumerModel)
}

func getAllConsumers(c *gin.Context) {
	consumerCollection := database.GetDatabase().Collection("consumers")

	cursor, err := consumerCollection.Find(c, bson.M{})

	if err != nil {
		panic("Failed to get consumer models: " + err.Error())
	}

	var consumerModels []model.Consumer
	err = cursor.All(c, &consumerModels)

	if err != nil {
		panic("Failed to get consumer models: " + err.Error())
	}

	if len(consumerModels) == 0 {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, consumerModels)
}

func getConsumerById(c *gin.Context) {
	consumerCollection := database.GetDatabase().Collection("consumers")

	id := c.Param("id")

	var consumerModel model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": id}).Decode(&consumerModel); err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, consumerModel)
}

func updateConsumerById(c *gin.Context) {
	var consumerRequestBody model.ConsumerRequestBody
	if err := c.ShouldBindJSON(&consumerRequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	consumerCollection := database.GetDatabase().Collection("consumers")

	id := c.Param("id")

	var consumerModel model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": id}).Decode(&consumerModel); err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	consumerModel.Name = consumerRequestBody.Name
	consumerModel.Location = consumerRequestBody.Location

	if _, err := consumerCollection.UpdateOne(c, bson.M{"_id": id}, bson.M{"$set": consumerModel}); err != nil {
		panic("Failed to update consumer model: " + err.Error())
	}

	c.JSON(http.StatusOK, consumerModel)
}

func deleteConsumerById(c *gin.Context) {
	consumerCollection := database.GetDatabase().Collection("consumers")

	id := c.Param("id")

	var consumerModel model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": id}).Decode(&consumerModel); err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if _, err := consumerCollection.DeleteOne(c, bson.M{"_id": id}); err != nil {
		panic("Failed to delete consumer model: " + err.Error())
	}

	c.JSON(http.StatusOK, consumerModel)
}

func AddExampleRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/", createConsumer)
	routerGroup.GET("/", getAllConsumers)
	routerGroup.GET(":id", getConsumerById)
	routerGroup.PUT(":id", updateConsumerById)
	routerGroup.DELETE(":id", deleteConsumerById)
}
