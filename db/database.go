package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect() {
	var uri string

	// Try to load .env file
	if err := godotenv.Load(); err == nil {
		uri = os.Getenv("MONGODB_URI")
	}

	// If uri is still empty, try to get it from environment variable
	if uri == "" {
		uri = os.Getenv("MONGODB_URI")
	}

	if uri == "" {
		log.Fatal("MongoDB URI is not set")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
}

func GetDatabase() *mongo.Database {
	return client.Database("watermark-generator")
}

func GetClient() *mongo.Client {
	return client
}

func Disconnect() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
