package activities

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ValidatePaymentRequest represents the input for payment validation
type ValidatePaymentRequest struct {
	PaymentID string  `json:"paymentId"`
	Amount    float64 `json:"amount"`
}

func (a *Activities) ValidatePayment(ctx context.Context, req ValidatePaymentRequest) (bool, error) {
	log.Printf("Validating payment: %s for amount: %.2f", req.PaymentID, req.Amount)
	
	// Simulate validation logic
	if req.PaymentID == "" {
		return false, fmt.Errorf("payment ID cannot be empty")
	}
	
	if req.Amount <= 0 {
		return false, fmt.Errorf("amount must be positive")
	}
	
	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)
	
	// In a real implementation, you'd validate payment methods, limits, etc.
	return true, nil
}