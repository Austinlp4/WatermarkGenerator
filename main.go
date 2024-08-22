package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"watermark-generator/api"
	"watermark-generator/db"
	"watermark-generator/watermark"

	"github.com/rs/cors"
)

//go:embed frontend/dist
var reactApp embed.FS

func init() {
	// Environment variables are already set by CapRover
	log.Println("Initializing application")
}

func main() {
	// Connect to MongoDB
	db.Connect()

	watermarkService := watermark.NewService()
	authHandler := api.NewAuthHandler()
	handler := api.NewWatermarkHandler(watermarkService)

	// Create a new mux for API routes
	apiMux := http.NewServeMux()
	api.SetupAuthRoutes(apiMux, authHandler) // Register auth routes with authHandler
	apiMux.HandleFunc("/api/watermark", handler.WatermarkHandler)
	apiMux.HandleFunc("/api/process-payment", handler.ProcessPaymentHandler)
	apiMux.HandleFunc("/api/create-checkout-session", handler.CreateCheckoutSessionHandler)
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
	})

	// Serve uploaded files
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080", "http://watermark-generator.com", "https://www.watermark-generator.com"},
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
