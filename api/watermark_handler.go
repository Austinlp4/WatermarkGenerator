package api

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"watermark-generator/db"
	"watermark-generator/watermark"

	"github.com/lucasb-eyer/go-colorful"
	"go.mongodb.org/mongo-driver/mongo"
)

type WatermarkHandler struct {
	service *watermark.Service
	logger  *log.Logger
	DB      *mongo.Database
}

func NewWatermarkHandler(service *watermark.Service) *WatermarkHandler {
	return &WatermarkHandler{
		service: service,
		logger:  log.New(os.Stdout, "API: ", log.LstdFlags),
		DB:      db.GetDatabase(),
	}
}

func (h *WatermarkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/watermark/text":
		h.TextWatermarkHandler(w, r)
	case "/api/watermark/image":
		h.ImageWatermarkHandler(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *WatermarkHandler) TextWatermarkHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("TextWatermarkHandler: Started processing request")
	defer h.logger.Println("TextWatermarkHandler: Finished processing request")

	if r.Method != http.MethodPost {
		h.logger.Println("TextWatermarkHandler: Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println("TextWatermarkHandler: Parsing multipart form")
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		h.logger.Printf("TextWatermarkHandler: Error parsing multipart form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	uniqueId := r.FormValue("uniqueId")
	if uniqueId == "" {
		h.logger.Println("TextWatermarkHandler: No uniqueId provided")
		http.Error(w, "No uniqueId provided", http.StatusBadRequest)
		return
	}

	h.logger.Printf("TextWatermarkHandler: Processing request with uniqueId: %s", uniqueId)

	// Log all form values
	h.logger.Println("Form values:")
	for key, values := range r.Form {
		h.logger.Printf("%s: %v", key, values)
	}

	h.logger.Println("TextWatermarkHandler: Retrieving file from form")
	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Printf("TextWatermarkHandler: Error retrieving file: %v", err)
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	h.logger.Printf("TextWatermarkHandler: File received: %s", header.Filename)

	opacity, err := strconv.ParseFloat(r.FormValue("opacity"), 64)
	if err != nil {
		h.logger.Printf("TextWatermarkHandler: Error parsing opacity: %v", err)
		opacity = 0.5 // Default opacity if parsing fails
	}
	opacity = math.Max(0, math.Min(1, opacity))

	spacing, err := strconv.ParseFloat(r.FormValue("spacing"), 64)
	if err != nil {
		h.logger.Printf("TextWatermarkHandler: Error parsing spacing: %v", err)
		spacing = 100 // Default spacing if parsing fails
	}

	text := r.FormValue("text")
	if text == "" {
		h.logger.Println("TextWatermarkHandler: No text provided for watermark")
		http.Error(w, "No text provided for watermark", http.StatusBadRequest)
		return
	}
	textColor := r.FormValue("color")
	if textColor == "" {
		textColor = "#000000" // Default to black if no color is provided
	}
	fontSize, err := strconv.ParseFloat(r.FormValue("fontSize"), 64)
	if err != nil {
		h.logger.Printf("TextWatermarkHandler: Error parsing fontSize: %v", err)
		fontSize = 32 // Default font size if parsing fails
	}

	h.logger.Println("TextWatermarkHandler: Calling ApplyWatermark")
	var result []byte
	result, err = h.service.ApplyWatermark(file, text, textColor, opacity, fontSize, spacing)

	if err != nil {
		h.logger.Printf("TextWatermarkHandler: Error applying watermark: %v", err)
		http.Error(w, fmt.Sprintf("Error applying watermark: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Printf("TextWatermarkHandler: Watermark applied successfully. Result length: %d", len(result))

	// Set cache control headers to prevent caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Set the content type and write the result
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("X-Unique-Id", uniqueId)
	w.Header().Set("Content-Length", strconv.Itoa(len(result)))

	h.logger.Println("TextWatermarkHandler: Headers set, attempting to write response")

	_, err = w.Write(result)
	if err != nil {
		h.logger.Printf("TextWatermarkHandler: Error writing response: %v", err)
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Println("TextWatermarkHandler: Response written successfully")
}

func (h *WatermarkHandler) ImageWatermarkHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("ImageWatermarkHandler: Started processing request")
	defer h.logger.Println("ImageWatermarkHandler: Finished processing request")

	if r.Method != http.MethodPost {
		h.logger.Println("ImageWatermarkHandler: Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println("ImageWatermarkHandler: Parsing multipart form")
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error parsing multipart form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	uniqueId := r.FormValue("uniqueId")
	if uniqueId == "" {
		h.logger.Println("ImageWatermarkHandler: No uniqueId provided")
		http.Error(w, "No uniqueId provided", http.StatusBadRequest)
		return
	}

	h.logger.Printf("ImageWatermarkHandler: Processing request with uniqueId: %s", uniqueId)

	// Log all form values
	h.logger.Println("Form values:")
	for key, values := range r.Form {
		h.logger.Printf("%s: %v", key, values)
	}

	h.logger.Println("ImageWatermarkHandler: Retrieving file from form")
	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error retrieving file: %v", err)
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	h.logger.Printf("ImageWatermarkHandler: File received: %s", header.Filename)

	opacity, err := strconv.ParseFloat(r.FormValue("opacity"), 64)
	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error parsing opacity: %v", err)
		opacity = 0.5 // Default opacity if parsing fails
	}
	opacity = math.Max(0, math.Min(1, opacity))

	spacing, err := strconv.ParseFloat(r.FormValue("spacing"), 64)
	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error parsing spacing: %v", err)
		spacing = 100 // Default spacing if parsing fails
	}

	watermarkSize, err := strconv.ParseFloat(r.FormValue("watermarkSize"), 64)
	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error parsing watermarkSize: %v", err)
		watermarkSize = 25 // Default watermark size if parsing fails
	}

	var result []byte
	watermarkImageFile, _, err := r.FormFile("watermarkImage")
	if err == nil {
		// Image watermark
		defer watermarkImageFile.Close()
		result, err = h.service.ApplyImageWatermark(file, watermarkImageFile, opacity, spacing, watermarkSize, uniqueId)
	} else {
		h.logger.Println("ImageWatermarkHandler: No watermark image provided")
		http.Error(w, "No watermark image provided", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error applying watermark: %v", err)
		http.Error(w, fmt.Sprintf("Error applying watermark: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Printf("ImageWatermarkHandler: Watermark applied successfully. Result length: %d", len(result))

	// Set cache control headers to prevent caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Set the content type and write the result
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("X-Unique-Id", uniqueId)
	w.Header().Set("Content-Length", strconv.Itoa(len(result)))

	h.logger.Println("ImageWatermarkHandler: Headers set, attempting to write response")

	_, err = w.Write(result)
	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error writing response: %v", err)
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Println("ImageWatermarkHandler: Response written successfully")
}

func (h *WatermarkHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	imagePath := r.URL.Query().Get("path")
	if imagePath == "" {
		http.Error(w, "Image path is required", http.StatusBadRequest)
		return
	}

	// Ensure the path is within the uploads directory
	fullPath := filepath.Join("./uploads", filepath.Clean(imagePath))
	if !strings.HasPrefix(fullPath, "./uploads/") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Serve the file
	http.ServeFile(w, r, fullPath)
}

func parseColor(s string) (color.Color, error) {
	c, err := colorful.Hex(s)
	if err != nil {
		return nil, err
	}
	return c, nil
}
