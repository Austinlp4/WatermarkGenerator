package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

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
	stripeHandler := api.NewStripeHandler(db.GetDatabase())

	// Create a new mux for API routes
	apiMux := http.NewServeMux()
	api.SetupAuthRoutes(apiMux, authHandler)
	apiMux.HandleFunc("/api/process-payment", handler.ProcessPaymentHandler)
	apiMux.HandleFunc("/api/create-checkout-session", handler.CreateCheckoutSessionHandler)
	apiMux.HandleFunc("/api/test-db", handler.TestDBConnectionHandler)
	apiMux.HandleFunc("/api/download", handler.DownloadHandler)
	apiMux.HandleFunc("/api/watermark/text", handler.TextWatermarkHandler)
	apiMux.HandleFunc("/api/watermark/image", handler.ImageWatermarkHandler)
	apiMux.HandleFunc("/api/create-subscription", stripeHandler.CreateSubscription)
	apiMux.HandleFunc("/api/cancel-subscription", stripeHandler.CancelSubscription)
	apiMux.HandleFunc("/api/webhook", stripeHandler.HandleWebhook)
	apiMux.HandleFunc("/api/watermark/bulk/text", handler.BulkTextWatermarkHandler)
	apiMux.HandleFunc("/api/watermark/bulk/image", handler.BulkImageWatermarkHandler)

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
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
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

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:8080",
		}, // Add your frontend URL here
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Wrap your main handler with the CORS handler
	corsHandler := c.Handler(mux)

	// Start the server
	log.Printf("Server starting on port 8080")
	http.ListenAndServe(":8080", corsHandler)
}
