package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"watermark-generator/models"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/subscription"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StripeHandler struct {
	DB *mongo.Database
}

func NewStripeHandler(db *mongo.Database) *StripeHandler {
	return &StripeHandler{DB: db}
}

func (h *StripeHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"userId"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received request with UserID: %s", req.UserID)

	// Fetch the user from the database
	objectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		log.Printf("Invalid UserID format: %v", err)
		http.Error(w, "Invalid UserID format", http.StatusBadRequest)
		return
	}

	var user models.User
	err = h.DB.Collection("users").FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		log.Printf("Error fetching user from database: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Set your Stripe secret key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Check if user has a Stripe Customer ID, if not, create one
	if user.StripeCustomerID == "" {
		params := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
		}
		customer, err := customer.New(params)
		if err != nil {
			log.Printf("Error creating Stripe customer: %v", err)
			http.Error(w, "Failed to create Stripe customer", http.StatusInternalServerError)
			return
		}
		user.StripeCustomerID = customer.ID

		// Update user in database with new Stripe Customer ID
		_, err = h.DB.Collection("users").UpdateOne(
			r.Context(),
			bson.M{"_id": objectID},
			bson.M{"$set": bson.M{"stripeCustomerId": customer.ID}},
		)
		if err != nil {
			log.Printf("Error updating user with Stripe Customer ID: %v", err)
			http.Error(w, "Failed to update user with Stripe Customer ID", http.StatusInternalServerError)
			return
		}
	}

	// Create a new Checkout Session for the subscription
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(user.StripeCustomerID),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(os.Getenv("PRO_PRICE_ID")), // Use PRO_PRICE_ID from env file
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(os.Getenv("VITE_API_URL") + "/subscribe/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(os.Getenv("VITE_API_URL") + "/subscribe/cancel"),
	}

	log.Printf("Creating Stripe session with params: %+v", params)

	session, err := session.New(params)
	if err != nil {
		log.Printf("Error creating Stripe session: %v", err)
		http.Error(w, fmt.Sprintf("Error creating Stripe session: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Stripe session created: %+v", session)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"sessionId": session.ID,
	}

	log.Printf("Sending response: %+v", response)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *StripeHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user from the database
	objectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = h.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Cancel the subscription in Stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	_, err = subscription.Cancel(user.SubscriptionId, &stripe.SubscriptionCancelParams{
		Prorate: stripe.Bool(false),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to cancel subscription in Stripe: %v", err), http.StatusInternalServerError)
		return
	}

	// Update user's subscription status in the database
	update := bson.M{
		"$set": bson.M{
			"subscriptionStatus":    "cancelled",
			"subscriptionExpiresAt": time.Now(),
			"subscriptionId":        "", // Clear the subscription ID
		},
	}
	_, err = h.DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update user subscription status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Subscription cancelled successfully"})
}

func (h *StripeHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Webhook handler called")

	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	log.Printf("Received webhook body: %s", string(body))

	var rawEvent struct {
		Type string `json:"type"`
		Data struct {
			Object json.RawMessage `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &rawEvent); err != nil {
		log.Printf("Error parsing webhook JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch rawEvent.Type {
	case "invoice.payment_succeeded":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(rawEvent.Data.Object, &session); err != nil {
			log.Printf("Error unmarshalling session data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Unmarshalled session data: %+v", session)
		err = h.handleSuccessfulSubscription(r.Context(), session)
		if err != nil {
			log.Printf("Error handling successful subscription: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Successfully handled subscription")
	default:
		log.Printf("Unhandled event type: %s", rawEvent.Type)
	}

	log.Println("Webhook handled successfully")
	w.WriteHeader(http.StatusOK)
}

func (h *StripeHandler) handleSuccessfulSubscription(ctx context.Context, session stripe.CheckoutSession) error {
	// Fetch subscription details
	subscription, err := subscription.Get(session.Subscription.ID, nil)
	if err != nil {
		return err
	}

	// Update user in database
	_, err = h.DB.Collection("users").UpdateOne(
		ctx,
		bson.M{"stripeCustomerId": session.Customer.ID},
		bson.M{"$set": bson.M{
			"subscriptionStatus":    "active",
			"subscriptionId":        subscription.ID,
			"subscriptionExpiresAt": time.Unix(subscription.CurrentPeriodEnd, 0),
		}},
	)
	return err
}
