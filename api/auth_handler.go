package api

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"log"
	"net/http"
	"strings"

	"watermark-generator/db"
	"watermark-generator/models"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *mongo.Database
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		DB: db.GetDatabase(),
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the entire request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Printf("Raw request body: %s", string(body))

	// Attempt to unmarshal the JSON
	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Log the received user data
	log.Printf("Received user data: %+v", user)

	// Check if password is empty
	if user.Password == "" {
		log.Println("Password is empty after JSON unmarshal")
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	// Trim the password
	user.Password = strings.TrimSpace(user.Password)
	log.Printf("Password after trim: '%s', length: %d", user.Password, len(user.Password))

	if user.Password == "" {
		log.Println("Password is empty after trimming")
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Update user object with hashed password
	user.Password = string(hashedPassword)

	// Generate a new ObjectID for the user
	user.ID = primitive.NewObjectID()

	// Insert the new user into the database
	collection := h.DB.Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Failed to insert user into database: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	log.Printf("User registered successfully: %s", user.Username)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Password length before trim: %d", len(credentials.Password))
	credentials.Password = strings.TrimSpace(credentials.Password)
	log.Printf("Password length after trim: %d", len(credentials.Password))

	log.Printf("Attempting to login user: %s", credentials.Username)
	log.Printf("Provided password: %s", credentials.Password)

	collection := h.DB.Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"username": credentials.Username}).Decode(&user)
	if err != nil {
		log.Printf("User not found: %s", credentials.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	log.Printf("Stored hashed password for user %s: %s", user.Username, user.Password)
	log.Printf("Provided password: %s", credentials.Password)

	// Check if the stored password looks like a bcrypt hash
	if !strings.HasPrefix(user.Password, "$2a$") && !strings.HasPrefix(user.Password, "$2b$") && !strings.HasPrefix(user.Password, "$2y$") {
		log.Printf("Warning: Stored password for user %s does not appear to be a bcrypt hash", user.Username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		log.Printf("Password comparison failed for user %s: %v", user.Username, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	log.Printf("Password match successful for user: %s", credentials.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func (h *AuthHandler) CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CurrentUserHandler called")

	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("No token provided")
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	// Remove 'Bearer ' prefix if present
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("Token received: %s", tokenString)

	// Validate the token
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("Token validated, subject: %s", claims.Subject)

	if claims.Subject == "" {
		log.Println("Empty subject in token")
		http.Error(w, "Invalid token: empty subject", http.StatusUnauthorized)
		return
	}

	// Use the subject as the user ID
	objectID, err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
		log.Printf("Invalid ObjectID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch user from database using claims.Subject (which should be the user ID)
	var user models.User
	err = h.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("User not found for ID: %s", claims.Subject)
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error fetching user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("User found: %s", user.Username)

	// Prepare the response
	response := map[string]string{
		"id":       user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Response sent successfully")
}

func (h *AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	credentials.Password = strings.TrimSpace(credentials.Password)

	collection := h.DB.Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"username": credentials.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a token that lasts for a week
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.Hex(), // Use "sub" instead of "user_id"
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Sign in successful",
		"username": user.Username,
		"token":    tokenString,
	})
}

func generateToken(userID string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 7 days
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *AuthHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch all users from the database
	collection := h.DB.Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err = cursor.All(context.Background(), &users); err != nil {
		http.Error(w, "Failed to decode users", http.StatusInternalServerError)
		return
	}

	// Return the users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *AuthHandler) DeleteAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := h.DB.Collection("users")
	result, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to delete users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "All users deleted successfully",
		"deletedCount": result.DeletedCount,
	})
}

func SetupAuthRoutes(mux *http.ServeMux, handler *AuthHandler) {
	mux.HandleFunc("/api/register", handler.RegisterHandler)
	mux.HandleFunc("/api/login", handler.LoginHandler)
	mux.HandleFunc("/api/signin", handler.SignInHandler)
	mux.HandleFunc("/api/current-user", handler.CurrentUserHandler)
	mux.HandleFunc("/api/users", handler.GetUsersHandler)
	mux.HandleFunc("/api/users/delete-all", handler.DeleteAllUsersHandler)
}
