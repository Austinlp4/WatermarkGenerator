package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"watermark-generator/api"
	"watermark-generator/db"
	"watermark-generator/watermark"

	"github.com/rs/cors"
)

//go:embed frontend/dist
var reactApp embed.FS

func main() {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/watermark-generator" // Default to localhost if not set
	}
	db.Connect(mongoURI)

	watermarkService := watermark.NewService()
	authHandler := &api.AuthHandler{DB: db.Client.Database("watermark-generator")} // Initialize AuthHandler
	handler := api.NewHandler(watermarkService, authHandler)

	// Create a new mux for API routes
	apiMux := http.NewServeMux()
	api.SetupAuthRoutes(apiMux, authHandler) // Register auth routes with authHandler
	apiMux.HandleFunc("/api/watermark", handler.WatermarkHandler.ServeHTTP)
	apiMux.HandleFunc("/api/process-payment", handler.ProcessPaymentHandler)
	apiMux.HandleFunc("/api/create-checkout-session", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateCheckoutSessionHandler(w, r)
	})
	apiMux.HandleFunc("/api/test-db", handler.TestDBConnectionHandler)

	// Create the main mux
	mux := http.NewServeMux()

	// Serve API routes
	mux.Handle("/api/", apiMux)

	// Serve React app
	fsys, err := fs.Sub(reactApp, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(fsys)))

	// Serve uploaded files
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // Enable debugging
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      c.Handler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Println("Server starting on :8080")
	log.Fatal(srv.ListenAndServe())
}
