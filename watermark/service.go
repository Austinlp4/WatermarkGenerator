package watermark

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"

	"database/sql"

	"github.com/golang/freetype/truetype"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"

	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Service struct {
	DB *sql.DB
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ApplyWatermark(r io.Reader, text string, textColor string, opacity float64, fontSize float64, spacing float64) ([]byte, error) {
	log.Printf("ApplyWatermark: Starting. Text: %s, Color: %s, Opacity: %.2f, Font Size: %.2f, Spacing: %.2f", text, textColor, opacity, fontSize, spacing)
	defer log.Println("ApplyWatermark: Finished")

	// Decode the original image
	srcImg, format, err := image.Decode(r)
	if err != nil {
		log.Printf("ApplyWatermark: Failed to decode source image: %v", err)
		return nil, fmt.Errorf("failed to decode source image: %v", err)
	}
	log.Printf("ApplyWatermark: Image decoded. Format: %s, Bounds: %v", format, srcImg.Bounds())

	result := image.NewRGBA(srcImg.Bounds())
	draw.Draw(result, srcImg.Bounds(), srcImg, image.Point{}, draw.Src)

	// Create and apply watermark
	if err := s.applyRepeatedWatermark(result, text, textColor, opacity, fontSize, spacing, 0, 0); err != nil {
		return nil, fmt.Errorf("failed to apply watermark: %v", err)
	}
	log.Printf("ApplyWatermark: Watermark applied to image")

	// Encode the result
	var buf bytes.Buffer
	if err := s.encodeImage(&buf, result, format); err != nil {
		return nil, fmt.Errorf("failed to encode result: %v", err)
	}
	log.Printf("ApplyWatermark: Image encoded. Buffer size: %d bytes", buf.Len())

	// Before returning
	log.Println("ApplyWatermark: Watermark applied successfully")
	return buf.Bytes(), nil
}

func (s *Service) applyRepeatedWatermark(img *image.RGBA, text string, textColor string, opacity float64, fontSize float64, spacing float64, width, height int) error {
	log.Printf("Applying repeated watermark. Text: %s, Opacity: %.2f, Font Size: %.2f, Spacing: %.2f", text, opacity, fontSize, spacing)

	bounds := img.Bounds()
	// Set base spacing appropriate for font size, then apply spacing multiplier
	baseSpacing := fontSize / 100
	verticalSpacing := int(baseSpacing * spacing)
	horizontalSpacing := int(baseSpacing * spacing)

	f, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return fmt.Errorf("failed to parse font: %v", err)
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size: fontSize,
		DPI:  72,
	})

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(applyOpacity(parseColor(textColor), opacity)),
		Face: face,
	}

	angle := 45.0

	for y := bounds.Min.Y - bounds.Max.Y; y < bounds.Max.Y*2; y += verticalSpacing {
		for x := bounds.Min.X - bounds.Max.X; x < bounds.Max.X*2; x += horizontalSpacing {
			err := drawRotatedText(d, text, x, y, angle)
			if err != nil {
				return fmt.Errorf("failed to draw rotated text: %v", err)
			}
		}
	}

	log.Printf("Watermark applied successfully")
	return nil
}

func drawRotatedText(d *font.Drawer, text string, x, y int, angle float64) error {
	// Convert angle to radians
	radians := angle * math.Pi / 180.0
	sin, cos := math.Sincos(radians)

	// Calculate the center point of the text
	width := d.MeasureString(text)
	centerX := float64(x) + float64(width.Round())/2
	centerY := float64(y)

	// Calculate rotated starting point
	startX := centerX*cos - centerY*sin
	startY := centerX*sin + centerY*cos

	// Set the starting point for the text
	d.Dot = fixed.Point26_6{
		X: fixed.I(int(startX)),
		Y: fixed.I(int(startY)),
	}

	// Draw the entire text string at once
	d.DrawString(text)

	return nil
}

func applyOpacity(c color.Color, opacity float64) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(float64(r>>8) * opacity),
		G: uint8(float64(g>>8) * opacity),
		B: uint8(float64(b>>8) * opacity),
		A: uint8(float64(a>>8) * opacity),
	}
}

func parseColor(colorStr string) color.Color {
	c, err := colorful.Hex(colorStr)
	if err != nil {
		return color.Black
	}
	return c
}

func (s *Service) encodeImage(w io.Writer, img image.Image, format string) error {
	log.Printf("Encoding image. Format: %s", format)
	switch format {
	case "jpeg":
		return jpeg.Encode(w, img, nil)
	case "png":
		return png.Encode(w, img)
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
}

func (s *Service) AuthenticateUser(email, password string) (*User, error) {
	var user User
	var hashedPassword string

	err := s.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func (s *Service) CreateSession(userID string) (string, error) {
	// Generate a random token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Set expiration time (e.g., 24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Store the session in the database
	_, err = s.DB.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		userID, token, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to store session: %v", err)
	}

	return token, nil
}

type User struct {
	ID    string
	Email string
}
