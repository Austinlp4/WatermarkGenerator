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
	// apiMux.HandleFunc("/api/watermark", handler.WatermarkHandler)
	apiMux.HandleFunc("/api/process-payment", handler.ProcessPaymentHandler)
	apiMux.HandleFunc("/api/create-checkout-session", handler.CreateCheckoutSessionHandler)
	apiMux.HandleFunc("/api/test-db", handler.TestDBConnectionHandler)
	apiMux.HandleFunc("/api/download", handler.DownloadHandler) // Add this line where you set up your routes
	apiMux.HandleFunc("/api/watermark/text", handler.TextWatermarkHandler)
	apiMux.HandleFunc("/api/watermark/image", handler.ImageWatermarkHandler)

	// Create the main mux
	mux := http.NewServeMux()

	// Serve API routes
	mux.Handle("/api/", apiMux)

	// Serve React app
	fsys, err := fs.Sub(reactApp, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(fsys))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiMux.ServeHTTP(w, r)
			return
		}
		// Serve index.html for all non-API routes
		if _, err := fsys.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			http.ServeFile(w, r, "frontend/dist/index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	// Serve uploaded files
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:8080",
			"http://watermark-generator.com",
			"https://www.watermark-generator.com",
			"https://watermark-generator.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // Enable debugging
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      c.Handler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second, // Increase this to 60 seconds
		IdleTimeout:  60 * time.Second,
	}
	log.Println("Server starting on :8080")
	log.Fatal(srv.ListenAndServe())
}
