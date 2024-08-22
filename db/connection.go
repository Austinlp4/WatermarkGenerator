package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

func TestConnection() error {
	// Implement your database connection test logic here
	// For example, using a global DB client:
	return Client.Ping(context.Background(), nil)
}
