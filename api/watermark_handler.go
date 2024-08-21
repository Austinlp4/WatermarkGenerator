package api

import (
	"image/color"
	"log"
	"net/http"
	"os"
	"strconv"

	"watermark-generator/models"
	"watermark-generator/watermark"

	"github.com/lucasb-eyer/go-colorful"
)

type WatermarkHandler struct {
	service *watermark.Service
	logger  *log.Logger
}

func NewWatermarkHandler(service *watermark.Service) *WatermarkHandler {
	return &WatermarkHandler{
		service: service,
		logger:  log.New(os.Stdout, "API: ", log.LstdFlags),
	}
}

func (h *WatermarkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.WatermarkHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WatermarkHandler) WatermarkHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("WatermarkHandler: Started processing request")
	defer h.logger.Println("WatermarkHandler: Finished processing request")

	if r.Method != http.MethodPost {
		h.logger.Println("WatermarkHandler: Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println("WatermarkHandler: Parsing multipart form")
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		h.logger.Printf("WatermarkHandler: Error parsing multipart form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Log all form values
	h.logger.Println("Form values:")
	for key, values := range r.Form {
		h.logger.Printf("%s: %v", key, values)
	}

	h.logger.Println("WatermarkHandler: Retrieving file from form")
	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Printf("WatermarkHandler: Error retrieving file: %v", err)
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	h.logger.Printf("WatermarkHandler: File received: %s", header.Filename)

	// Extract other form values
	text := r.FormValue("text")
	textColor := r.FormValue("color")
	h.logger.Printf("WatermarkHandler: Received color: %s", textColor)
	if textColor == "" {
		textColor = "#000000" // Default to black if no color is provided
		h.logger.Println("WatermarkHandler: Using default color #000000")
	}

	opacity, err := strconv.ParseFloat(r.FormValue("opacity"), 64)
	if err != nil {
		h.logger.Printf("WatermarkHandler: Error parsing opacity: %v", err)
		opacity = 0.5 // Default opacity if parsing fails
	}

	fontSize, err := strconv.ParseFloat(r.FormValue("fontSize"), 64)
	if err != nil {
		h.logger.Printf("WatermarkHandler: Error parsing fontSize: %v", err)
		fontSize = 32 // Default font size if parsing fails
	}

	spacing, err := strconv.ParseFloat(r.FormValue("spacing"), 64)
	if err != nil {
		h.logger.Printf("WatermarkHandler: Error parsing spacing: %v", err)
		spacing = 100 // Default spacing if parsing fails
	}

	h.logger.Printf("WatermarkHandler: Parsed values - Text: %s, Color: %s, Opacity: %.2f, Font Size: %.2f, Spacing: %.2f", text, textColor, opacity, fontSize, spacing)

	// Parse the text color
	color, err := parseColor(textColor)
	if err != nil {
		h.logger.Printf("WatermarkHandler: Error parsing color: %v", err)
		http.Error(w, "Invalid color format", http.StatusBadRequest)
		return
	}

	h.logger.Printf("WatermarkHandler: Applying watermark. Text: %s, Color: %s, Opacity: %.2f, Font Size: %.2f, Spacing: %.2f", text, textColor, opacity, fontSize, spacing)

	// Call the watermark service
	colorStr := color.(colorful.Color).Hex() // Cast to colorful.Color
	result, err := h.service.ApplyWatermark(file, text, colorStr, opacity, fontSize, spacing)
	if err != nil {
		h.logger.Printf("WatermarkHandler: Error applying watermark: %v", err)
		http.Error(w, "Error applying watermark", http.StatusInternalServerError)
		return
	}

	h.logger.Println("WatermarkHandler: Watermark applied successfully")

	// Set the content type and write the result
	w.Header().Set("Content-Type", "image/png")
	w.Write(result)
}

func (h *WatermarkHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context
	user := r.Context().Value("user").(*models.User)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the image path from the query parameters
	imagePath := r.URL.Query().Get("path")
	if imagePath == "" {
		http.Error(w, "Image path is required", http.StatusBadRequest)
		return
	}

	// Serve the file
	http.ServeFile(w, r, imagePath)
}

func parseColor(s string) (color.Color, error) {
	c, err := colorful.Hex(s)
	if err != nil {
		return nil, err
	}
	return c, nil
}
