package routes

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/segmentio/kafka-go"
)

func getKafkaReader(topic string) *kafka.Reader {
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	kafkaURL := fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)

	if kafkaHost == "" || kafkaPort == "" {
		kafkaURL = "localhost:9092"
	}

	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		GroupID:   "products-consumer-group",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
}

func ProductsConsumer() {
	reader := getKafkaReader("products")
	defer reader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		// call the function to update the product
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
