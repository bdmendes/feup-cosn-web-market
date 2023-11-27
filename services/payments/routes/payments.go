package routes

import (
	"cosn/payments/model"
	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

func payOrder(c *gin.Context) {
	var paymentRequestData model.PaymentRequestData
	if err := c.ShouldBindJSON(&paymentRequestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if model.IsPaymentMethodValid(paymentRequestData.PaymentMethod) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment method not valid"})
		return
	}

	// Currently, the service does not work with real data.
	// It fails 20% of the time.
	r, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		panic(err)
	}
	if r.Cmp(big.NewInt(20)) < 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment failed. try again later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "payment successful"})
}

func AddPaymentRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/", payOrder)
}
