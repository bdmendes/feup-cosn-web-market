package routes

import (
	"context"
	"cosn/orders/model"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

func getKafkaWriter() *kafka.Writer {
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	kafkaURL := fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)

	if kafkaHost == "" || kafkaPort == "" {
		kafkaURL = "localhost:9092"
	}

	topic := os.Getenv("KAFKA_ORDER_NOTIFICATIONS_TOPIC")
	if topic == "" {
		topic = "notify-order"
	}

	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func PublishOrder(order model.Order) {
	writer := getKafkaWriter()

	productNotificationBytes, err := json.Marshal(order.Products)
	if err != nil {
		fmt.Printf("error marshalling product notification: %v\n", err)
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: productNotificationBytes,
	})
	if err != nil {
		fmt.Printf("error publishing product: %v\n", err)
		return
	}
}
