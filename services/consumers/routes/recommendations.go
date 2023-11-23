package routes

import (
	"cosn/consumers/database"
	"cosn/consumers/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getRecommendedProducts(c *gin.Context, numberOfRecommendations int) {
	consumerId := c.Param("consumerId")

	consumerCollection := database.GetDatabase().Collection("consumers")

	var consumer model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": consumerId}).Decode(&consumer); err != nil {
		panic("Failed to get consumer: " + err.Error())
	}

	relatedProducts := consumer.RelatedProducts(c)
	var recommendations []model.Product

	for _, product := range relatedProducts {
		if numberOfRecommendations > 0 {
			recommendations = append(recommendations, product)
			numberOfRecommendations--
		}
	}

	c.JSON(http.StatusOK, recommendations)
}

func getRecommendationsForCategory(c *gin.Context, numberOfRecommendations int) {
	consumerId := c.Param("consumerId")
	category := c.Param("category")

	consumerCollection := database.GetDatabase().Collection("consumers")

	var consumer model.Consumer
	if err := consumerCollection.FindOne(c, bson.M{"_id": consumerId}).Decode(&consumer); err != nil {
		panic("Failed to get consumer: " + err.Error())
	}

	relatedProducts := consumer.RelatedProducts(c)
	var relatedProductsOfCategory []model.Product
	for _, product := range relatedProducts {
		if product.Category == category && numberOfRecommendations > 0 {
			relatedProductsOfCategory = append(relatedProductsOfCategory, product)
			numberOfRecommendations--
		}
	}

	c.JSON(http.StatusOK, relatedProductsOfCategory)
}

func AddRecommendationRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET(":consumerId/recommendations/category/:category", func(c *gin.Context) {
		numberOfRecommendations := c.Query("size")
		if numberOfRecommendations == "" {
			getRecommendationsForCategory(c, 10)
		} else {
			numberOfRecommendationsInt, err := strconv.ParseInt(numberOfRecommendations, 10, 32)
			if err != nil {
				panic("Failed to parse number of recommendations: " + err.Error())
			}
			getRecommendationsForCategory(c, int(numberOfRecommendationsInt))
		}
	})
	routerGroup.GET(":consumerId/recommendations", func(c *gin.Context) {
		numberOfRecommendations := c.Query("size")
		if numberOfRecommendations == "" {
			getRecommendedProducts(c, 10)
		} else {
			numberOfRecommendationsInt, err := strconv.ParseInt(numberOfRecommendations, 10, 32)
			if err != nil {
				panic("Failed to parse number of recommendations: " + err.Error())
			}
			getRecommendedProducts(c, int(numberOfRecommendationsInt))
		}
	})
}
