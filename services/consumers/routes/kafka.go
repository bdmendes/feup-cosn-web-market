package routes

import (
	"context"
	"cosn/consumers/database"
	"cosn/consumers/model" // #nosec G501
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
)

func getKafkaReader() *kafka.Reader {
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	kafkaURL := fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)

	if kafkaHost == "" || kafkaPort == "" {
		kafkaURL = "localhost:9092"
	}

	topic := os.Getenv("KAFKA_PRODUCTS_TOPIC")
	if topic == "" {
		topic = "purchasedproducts"
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
	reader := getKafkaReader()
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("error reading message: %v\n", err)
		}

		var productNotification model.ProductNotification
		if err = json.Unmarshal(m.Value, &productNotification); err != nil {
			fmt.Printf("error unmarshalling product notification: %v\n", err)
		}

		createOrUpdateProduct(productNotification)
	}
}

func createOrUpdateProduct(productNotification model.ProductNotification) {
	productsCollection := database.GetDatabase().Collection("products")

	var product model.Product
	err := productsCollection.FindOne(context.Background(), bson.M{"id": productNotification.ID}).Decode(&product)
	if err != nil { // create product
		product.ID = productNotification.ID
		product.Category = productNotification.Category
		product.Name = productNotification.Name
		product.Brand = productNotification.Brand
		product.Prices = []float32{productNotification.Price}

		if _, err = productsCollection.InsertOne(context.Background(), product); err != nil {
			fmt.Printf("error inserting product: %v\n", err)
		}
	} else { // update product
		product.Category = productNotification.Category
		product.Name = productNotification.Name
		product.Brand = productNotification.Brand
		product.Prices = append(product.Prices, productNotification.Price)

		if _, err = productsCollection.UpdateOne(context.Background(),
			bson.M{"id": productNotification.ID}, bson.M{"$set": product}); err != nil {
			fmt.Printf("error updating product: %v\n", err)
		}
	}
}

// How to run Kafka:
// bin/kafka-server-start.sh config/server.properties
// bin/zookeeper-server-start.sh config/zookeeper.properties
// bin/kafka-topics.sh --create --topic products --bootstrap-server localhost:9092
// bin/kafka-console-producer.sh --topic products --bootstrap-server localhost:9092
