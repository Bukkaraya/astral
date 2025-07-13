package activities

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// ChargePaymentRequest represents the input for payment charging
type ChargePaymentRequest struct {
	PaymentID string  `json:"paymentId"`
	Amount    float64 `json:"amount"`
}

func (a *Activities) ChargePayment(ctx context.Context, req ChargePaymentRequest) (string, error) {
	log.Printf("Charging payment: %s for amount: %.2f", req.PaymentID, req.Amount)
	
	// Simulate payment processing
	time.Sleep(500 * time.Millisecond)
	
	transactionID := fmt.Sprintf("txn_%s_%s", req.PaymentID, uuid.New().String()[:8])
	
	// In a real implementation, you'd integrate with payment gateways
	log.Printf("Payment charged successfully. Transaction ID: %s", transactionID)
	
	return transactionID, nil
}