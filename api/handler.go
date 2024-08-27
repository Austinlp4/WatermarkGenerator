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
	WatermarkHandler *WatermarkHandler
	AuthHandler      *AuthHandler
	mux              *http.ServeMux
}

func NewHandler(service *watermark.Service, authHandler *AuthHandler) *Handler {
	h := &Handler{
		WatermarkHandler: NewWatermarkHandler(service),
		AuthHandler:      authHandler,
		mux:              http.NewServeMux(),
	}
	h.setupRoutes()
	return h
}

func (h *Handler) setupRoutes() {
	h.mux.HandleFunc("/api/watermark/text", h.WatermarkHandler.TextWatermarkHandler)
	h.mux.HandleFunc("/api/watermark/image", h.WatermarkHandler.ImageWatermarkHandler)
	// Add other routes as needed
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in ServeHTTP: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	h.mux.ServeHTTP(w, r)
	log.Printf("Finished processing request: %s %s", r.Method, r.URL.Path)
}

func SetupRoutes(mux *http.ServeMux, handler *Handler) http.Handler {
	mux.HandleFunc("/api/signin", handler.AuthHandler.LoginHandler)
	mux.HandleFunc("/api/watermark/image", LoggingMiddleware(AuthMiddleware(handler.WatermarkHandler.ImageWatermarkHandler)))
	mux.HandleFunc("/api/watermark/text", LoggingMiddleware(AuthMiddleware(handler.WatermarkHandler.TextWatermarkHandler)))
	mux.HandleFunc("/api/download", LoggingMiddleware(AuthMiddleware(handler.WatermarkHandler.DownloadHandler)))

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

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("LoggingMiddleware: Received request for %s", r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("LoggingMiddleware: Finished processing request for %s", r.URL.Path)
	}
}
