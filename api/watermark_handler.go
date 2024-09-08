package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"watermark-generator/db"
	"watermark-generator/models"
	"watermark-generator/watermark"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	case "/api/watermark/bulk/text":
		h.BulkTextWatermarkHandler(w, r)
	case "/api/watermark/bulk/image":
		h.BulkImageWatermarkHandler(w, r)
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

	var uniqueId string
	if ctxUniqueId, ok := r.Context().Value("uniqueId").(string); ok && ctxUniqueId != "" {
		uniqueId = ctxUniqueId
	} else {
		uniqueId = r.FormValue("uniqueId")
	}

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

	// Set the content type and other headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Unique-Id", uniqueId)

	h.logger.Println("TextWatermarkHandler: Headers set, attempting to write response")

	// Create a response similar to the bulk handler
	response := map[string]interface{}{
		"message": "Watermark applied successfully",
		"results": []map[string]interface{}{
			{
				"filename": header.Filename,
				"data":     fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(result)),
				"uniqueId": uniqueId,
			},
		},
	}

	// Encode and write the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("TextWatermarkHandler: Error encoding JSON response: %v", err)
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
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

	var uniqueId string
	if ctxUniqueId, ok := r.Context().Value("uniqueId").(string); ok && ctxUniqueId != "" {
		uniqueId = ctxUniqueId
	} else {
		uniqueId = r.FormValue("uniqueId")
	}

	if uniqueId == "" {
		h.logger.Println("ImageWatermarkHandler: No uniqueId provided")
		http.Error(w, "No uniqueId provided", http.StatusBadRequest)
		return
	}

	h.logger.Printf("ImageWatermarkHandler: Processing request with uniqueId: %s", uniqueId)

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

	watermarkImageFile, watermarkHeader, err := r.FormFile("watermarkImage")
	if err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error retrieving watermark image: %v", err)
		http.Error(w, "No watermark image provided", http.StatusBadRequest)
		return
	}
	defer watermarkImageFile.Close()

	h.logger.Printf("ImageWatermarkHandler: Watermark image received: %s", watermarkHeader.Filename)

	h.logger.Println("ImageWatermarkHandler: Calling ApplyImageWatermark")
	result, err := h.service.ApplyImageWatermark(file, watermarkImageFile, opacity, spacing, watermarkSize, uniqueId)
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

	// Set the content type and other headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Unique-Id", uniqueId)

	h.logger.Println("ImageWatermarkHandler: Headers set, attempting to write response")

	// Create a response similar to the text watermark handler
	response := map[string]interface{}{
		"message": "Watermark applied successfully",
		"results": []map[string]interface{}{
			{
				"filename": header.Filename,
				"data":     fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(result)),
				"uniqueId": uniqueId,
			},
		},
	}

	// Encode and write the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("ImageWatermarkHandler: Error encoding JSON response: %v", err)
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Println("ImageWatermarkHandler: Response written successfully")
}

func (h *WatermarkHandler) BulkTextWatermarkHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("BulkTextWatermarkHandler: Started processing request")
	defer h.logger.Println("BulkTextWatermarkHandler: Finished processing request")

	if r.Method != http.MethodPost {
		h.logger.Println("BulkTextWatermarkHandler: Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println("BulkTextWatermarkHandler: Parsing multipart form")
	err := r.ParseMultipartForm(50 << 20) // 50 MB
	if err != nil {
		h.logger.Printf("BulkTextWatermarkHandler: Error parsing multipart form: %v", err)
		http.Error(w, fmt.Sprintf("Unable to parse form: %v", err), http.StatusBadRequest)
		return
	}

	userId := r.FormValue("userId")
	if userId == "" {
		h.logger.Println("BulkTextWatermarkHandler: No userId provided")
		http.Error(w, "No userId provided", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		h.logger.Println("BulkTextWatermarkHandler: No files provided")
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	var results []map[string]interface{}
	for _, fileHeader := range files {
		uniqueId := uuid.New().String()

		// Create a new context with the uniqueId
		ctx := context.WithValue(r.Context(), "uniqueId", uniqueId)

		// Create a new request with the updated context
		fileRequest := r.WithContext(ctx)
		fileRequest.MultipartForm = &multipart.Form{
			File: map[string][]*multipart.FileHeader{
				"image": {fileHeader},
			},
			Value: r.MultipartForm.Value,
		}
		// Add uniqueId to form values
		fileRequest.MultipartForm.Value["uniqueId"] = []string{uniqueId}

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// Call the TextWatermarkHandler
		h.TextWatermarkHandler(rr, fileRequest)

		// Check if the watermarking was successful
		if rr.Code != http.StatusOK {
			h.logger.Printf("BulkTextWatermarkHandler: Error processing file %s: %s", fileHeader.Filename, rr.Body.String())
			continue
		}

		// Parse the JSON response
		var response map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			h.logger.Printf("BulkTextWatermarkHandler: Error parsing response for file %s: %v", fileHeader.Filename, err)
			continue
		}

		// Extract the result from the response
		if resultArray, ok := response["results"].([]interface{}); ok && len(resultArray) > 0 {
			if result, ok := resultArray[0].(map[string]interface{}); ok {
				results = append(results, result)
			}
		}
	}

	h.logger.Println("BulkTextWatermarkHandler: Watermark applied successfully to all images")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Watermark applied successfully",
		"results": results,
	})
}

func (h *WatermarkHandler) BulkImageWatermarkHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("BulkImageWatermarkHandler: Started processing request")
	defer h.logger.Println("BulkImageWatermarkHandler: Finished processing request")

	if r.Method != http.MethodPost {
		h.logger.Println("BulkImageWatermarkHandler: Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println("BulkImageWatermarkHandler: Parsing multipart form")
	err := r.ParseMultipartForm(50 << 20) // 50 MB
	if err != nil {
		h.logger.Printf("BulkImageWatermarkHandler: Error parsing multipart form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	userId := r.FormValue("userId")
	if userId == "" {
		h.logger.Println("BulkImageWatermarkHandler: No userId provided")
		http.Error(w, "No userId provided", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		h.logger.Println("BulkImageWatermarkHandler: No files provided")
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	// Get the watermark image once
	watermarkImageFile, watermarkImageHeader, err := r.FormFile("watermarkImage")
	if err != nil {
		h.logger.Println("BulkImageWatermarkHandler: No watermark image provided")
		http.Error(w, "No watermark image provided", http.StatusBadRequest)
		return
	}
	defer watermarkImageFile.Close()

	var results []map[string]interface{}
	for _, fileHeader := range files {
		uniqueId := uuid.New().String()

		// Create a new context with the uniqueId
		ctx := context.WithValue(r.Context(), "uniqueId", uniqueId)

		// Create a new request with the updated context
		fileRequest := r.WithContext(ctx)
		fileRequest.MultipartForm = &multipart.Form{
			File: map[string][]*multipart.FileHeader{
				"image":          {fileHeader},
				"watermarkImage": {watermarkImageHeader},
			},
			Value: r.MultipartForm.Value,
		}
		// Add uniqueId to form values
		fileRequest.MultipartForm.Value["uniqueId"] = []string{uniqueId}

		// Create a ResponseRecorder to capture the response
		rr := httptest.NewRecorder()

		// Call the ImageWatermarkHandler
		h.ImageWatermarkHandler(rr, fileRequest)

		// Check if the watermarking was successful
		if rr.Code != http.StatusOK {
			h.logger.Printf("BulkImageWatermarkHandler: Error processing file %s: %s", fileHeader.Filename, rr.Body.String())
			continue
		}

		// Parse the JSON response
		var response map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			h.logger.Printf("BulkImageWatermarkHandler: Error parsing response for file %s: %v", fileHeader.Filename, err)
			continue
		}

		// Extract the result from the response
		if resultArray, ok := response["results"].([]interface{}); ok && len(resultArray) > 0 {
			if result, ok := resultArray[0].(map[string]interface{}); ok {
				results = append(results, result)
			}
		}
	}

	h.logger.Println("BulkImageWatermarkHandler: Watermark applied successfully to all images")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Watermark applied successfully",
		"results": results,
	})
}

func (h *WatermarkHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from the request context or token
	userID := r.Context().Value("userID").(string)
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Fetch the user from the database
	var user models.User
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = h.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check subscription status and daily download limit
	if user.SubscriptionStatus != "active" || user.SubscriptionExpiresAt.Before(time.Now()) {
		today := time.Now().Truncate(24 * time.Hour)
		if user.LastDownloadDate.Before(today) {
			// Reset daily downloads if it's a new day
			user.DailyDownloads = 0
		}
		if user.DailyDownloads >= 1 {
			http.Error(w, "Daily download limit reached", http.StatusForbidden)
			return
		}
		user.DailyDownloads++
		user.LastDownloadDate = time.Now()
		_, err = h.DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{
			"$set": bson.M{
				"dailyDownloads":   user.DailyDownloads,
				"lastDownloadDate": user.LastDownloadDate,
			},
		})
		if err != nil {
			http.Error(w, "Failed to update user download count", http.StatusInternalServerError)
			return
		}
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
