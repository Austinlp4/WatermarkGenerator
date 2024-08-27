package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"watermark-generator/db"
)

type User struct {
	ID       string
	Username string
}

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Checking authorization header")
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Println("AuthMiddleware: No authorization header found")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Remove 'Bearer ' prefix if present
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validate the token
		user, err := validateToken(tokenString)
		if err != nil {
			log.Printf("AuthMiddleware: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set the user in the request context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		log.Println("AuthMiddleware: User authenticated successfully")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func validateToken(tokenString string) (*User, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Validate the token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user ID from the token claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		// Fetch the user from the database
		user, err := GetUserByID(userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %v", err)
		}

		return user, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func GetUserByID(userID string) (*User, error) {
	collection := db.GetDatabase().Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with ID: %s", userID)
		}
		return nil, fmt.Errorf("error querying user: %v", err)
	}

	return &user, nil
}
