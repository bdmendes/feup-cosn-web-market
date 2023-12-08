package routes

import (
	"context"
	"cosn/consumers/database"
	"cosn/consumers/model"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
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

	hasher := md5.New()
	hasher.Write([]byte(productNotification.ID))
	hash := hasher.Sum(nil)

	product_id, err := primitive.ObjectIDFromHex(hex.EncodeToString(hash)[:24])
	if err != nil {
		fmt.Printf("error converting product id: %v\n", err)
		return
	}

	var product model.Product
	err = productsCollection.FindOne(context.Background(), bson.M{"_id": product_id}).Decode(&product)
	if err != nil { // create product
		product.ID = product_id
		product.Category = productNotification.Category
		product.Description = productNotification.Description
		product.Prices = []float32{productNotification.Price}

		if _, err = productsCollection.InsertOne(context.Background(), product); err != nil {
			fmt.Printf("error inserting product: %v\n", err)
		}
	} else { // update product
		product.Category = productNotification.Category
		product.Description = productNotification.Description
		product.Prices = append(product.Prices, productNotification.Price)

		if _, err = productsCollection.UpdateOne(context.Background(), bson.M{"_id": product_id}, bson.M{"$set": product}); err != nil {
			fmt.Printf("error updating product: %v\n", err)
		}
	}
}

// How to run Kafka:
// bin/kafka-server-start.sh config/server.properties
// bin/zookeeper-server-start.sh config/zookeeper.properties
// bin/kafka-topics.sh --create --topic products --bootstrap-server localhost:9092
// bin/kafka-console-producer.sh --topic products --bootstrap-server localhost:9092
