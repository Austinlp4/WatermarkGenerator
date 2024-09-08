package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

// Add this method to the WatermarkHandler struct
func (h *WatermarkHandler) ProcessPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var paymentRequest struct {
		PaymentMethodID string  `json:"paymentMethodId"`
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
		Description     string  `json:"description"`
	}

	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Load Stripe secret key from environment variable
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(paymentRequest.Amount)),
		Currency:      stripe.String(paymentRequest.Currency),
		PaymentMethod: stripe.String(paymentRequest.PaymentMethodID),
		Description:   stripe.String(paymentRequest.Description),
		Confirm:       stripe.Bool(true),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating payment intent: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"paymentIntentID": pi.ID,
		"clientSecret":    pi.ClientSecret,
	})
}

func (h *WatermarkHandler) CreateCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(req.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Donation to Watermark Wizard"),
					},
					UnitAmount: stripe.Int64(req.Amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(os.Getenv("VITE_API_URL") + "/?donation=success"),
		CancelURL:  stripe.String(os.Getenv("VITE_API_URL") + "/?donation=cancelled"),
	}

	session, err := session.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"sessionId": session.ID,
	})
}
