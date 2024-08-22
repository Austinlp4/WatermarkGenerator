package api

import (
	"context"
	"watermark-generator/db"
	"watermark-generator/watermark"

	"log"

	"net/http"

	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	db.Connect()

	client := db.GetClient()
	if client == nil {
		log.Fatal("Failed to connect to MongoDB")
	}

	log.Println("Successfully connected to MongoDB")

	// Test the connection by performing a simple operation
	collection := client.Database("watermark-generator").Collection("test")
	_, err := collection.InsertOne(context.Background(), bson.M{"test": "connection"})
	if err != nil {
		log.Fatal("Failed to insert test document:", err)
	}

	log.Println("Successfully inserted test document")
}

type Handler struct {
	*WatermarkHandler
	*AuthHandler
}

func NewHandler(service *watermark.Service, authHandler *AuthHandler) *Handler {
	return &Handler{
		WatermarkHandler: NewWatermarkHandler(service),
		AuthHandler:      authHandler,
	}
}

func SetupRoutes(mux *http.ServeMux, handler *Handler) http.Handler {
	mux.HandleFunc("/api/signin", handler.SignInHandler)
	mux.HandleFunc("/api/watermark", AuthMiddleware(handler.WatermarkHandler.WatermarkHandler))
	mux.HandleFunc("/api/download", AuthMiddleware(handler.WatermarkHandler.DownloadHandler))

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	return c.Handler(mux)
}
