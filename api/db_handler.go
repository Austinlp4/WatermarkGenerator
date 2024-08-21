package api

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) TestDBConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := db.Database("watermark-generator").Collection("test_collection")

	// Insert a test document
	testDoc := bson.M{"name": "test", "value": "This is a test document"}
	_, err := collection.InsertOne(context.Background(), testDoc)
	if err != nil {
		http.Error(w, "Failed to insert document", http.StatusInternalServerError)
		return
	}

	// Retrieve the test document
	var result bson.M
	err = collection.FindOne(context.Background(), bson.M{"name": "test"}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No document found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve document", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the retrieved document
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
