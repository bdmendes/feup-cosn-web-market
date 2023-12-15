package tasks

import (
	"context"
	"cosn/orders/database"
	"cosn/orders/model"
	"cosn/orders/routes"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ScheduledOrdersTask() {
	ordersCollection := database.GetDatabase().Collection("orders")

	for {
		// For production:
		// time.Sleep(12 * time.Hour)

		// For testing:
		time.Sleep(15 * time.Second)

		cursor, err := ordersCollection.Find(context.Background(), bson.M{"intervaldays": bson.M{"$gt": 0}})
		if err != nil {
			println("Order Task: " + err.Error())
			continue
		}

		var orders []model.Order
		err = cursor.All(context.Background(), &orders)
		if err != nil {
			println("Order Task: " + err.Error())
			continue
		}

		for _, order := range orders {
			if time.Now().After(order.Date.AddDate(0, 0, order.IntervalDays)) {
				new_order := order
				new_order.ID = primitive.NewObjectID()
				new_order.Status = model.PENDING
				new_order.Date = time.Now()

				order.IntervalDays = 0
				update := bson.M{
					"$set": bson.M{
						"intervaldays": order.IntervalDays,
					},
				}

				doc := ordersCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": order.ID}, update)
				if doc.Err() != nil {
					println("Order Task: " + doc.Err().Error())
					break
				}

				if _, err := ordersCollection.InsertOne(context.Background(), new_order); err != nil {
					println("Order Task: " + err.Error())
					break
				}

				payload := []byte(`{"amount": ` + strconv.FormatFloat(new_order.PaymentData.Amount, 'f', -1, 64) +
					`, "payment_method": "` + new_order.PaymentData.PaymentMethod +
					`", "payment_data": "` + new_order.PaymentData.PaymentData +
					`"}`)

				go routes.SendPostRequest(payload, os.Getenv("PAYMENT_SERVICE_URL")+"/payment", routes.PaymentCallback, map[string]interface{}{"order_id": new_order.ID.Hex()})
			}
		}
	}
}
