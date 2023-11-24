package tasks

import (
	"context"
	"cosn/orders/database"
	"cosn/orders/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ScheduledOrdersTask() {
	ordersCollection := database.GetDatabase().Collection("orders")

	for {
		time.Sleep(1 * time.Second)

		cursor, err := ordersCollection.Find(context.Background(), bson.M{"intervaldays": bson.M{"$gt": 0}})
		if err != nil {
			continue
		}

		var orders []model.Order
		err = cursor.All(context.Background(), &orders)
		if err != nil {
			continue
		}

		for _, order := range orders {
			if time.Now().After(order.Date.AddDate(0, 0, *order.IntervalDays)) {
				new_order := order

				new_order.ID = primitive.NewObjectID()

				new_status := model.PENDING
				new_order.Status = &new_status

				new_date := time.Now()
				new_order.Date = &(new_date)

				zero := 0
				order.IntervalDays = &zero
				update := bson.M{
					"$set": bson.M{
						"intervaldays": order.IntervalDays,
					},
				}

				doc := ordersCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": order.ID}, update)
				if doc.Err() != nil {
					break
				}

				if _, err := ordersCollection.InsertOne(context.Background(), new_order); err != nil {
					break
				}
			}
		}
	}
}
