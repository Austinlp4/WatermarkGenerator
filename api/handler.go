package api

import (
	"context"
	"os"
	"watermark-generator/watermark"

	"log"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize MongoDB client with authentication
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetAuth(options.Credential{
		Username:      "root",
		Password:      "Seahawks",
		AuthMechanism: "SCRAM-SHA-256",
	})
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client
}

type Handler struct {
	*WatermarkHandler
	*AuthHandler // Add AuthHandler to the Handler struct
}

func NewHandler(service *watermark.Service, authHandler *AuthHandler) *Handler {
	return &Handler{
		WatermarkHandler: NewWatermarkHandler(service),
		AuthHandler:      authHandler, // Initialize AuthHandler
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
