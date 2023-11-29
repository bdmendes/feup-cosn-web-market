package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func GetDatabase() *mongo.Database {
	if database == nil {
		panic("Database not initialized")
	}

	return database
}

func InitDatabase() *mongo.Database {
	databaseUrl := os.Getenv("MONGO_URL")
	if databaseUrl == "" {
		panic("Missing MONGO_URL env variable")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(databaseUrl))
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic("Failed to ping MongoDB: " + err.Error())
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		panic("Missing MONGO_DATABASE env variable")
	}

	database = client.Database(databaseName)
	if database == nil {
		panic("Failed to get database")
	}

	return database
}

func DisconnectDatabase() {
	err := database.Client().Disconnect(context.Background())
	if err != nil {
		panic("Failed to disconnect from MongoDB: " + err.Error())
	}
}
