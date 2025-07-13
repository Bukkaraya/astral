package activities

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// ProcessRefundRequest represents the input for refund processing
type ProcessRefundRequest struct {
	PaymentID string `json:"paymentId"`
}

func (a *Activities) ProcessRefund(ctx context.Context, req ProcessRefundRequest) (string, error) {
	log.Printf("Processing refund for payment: %s", req.PaymentID)
	
	// Simulate refund processing
	time.Sleep(400 * time.Millisecond)
	
	refundID := fmt.Sprintf("ref_%s_%s", req.PaymentID, uuid.New().String()[:8])
	
	// In a real implementation, you'd process refunds through payment gateways
	log.Printf("Refund processed successfully. Refund ID: %s", refundID)
	
	return refundID, nil
}