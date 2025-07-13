package activities

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type ProcessShippingRequest struct {
	OrderID string `json:"orderId"`
}

func (a *Activities) ProcessShipping(ctx context.Context, req ProcessShippingRequest) (string, error) {
	log.Printf("Processing shipping for order: %s", req.OrderID)
	
	// Simulate shipping processing
	time.Sleep(300 * time.Millisecond)
	
	shippingID := fmt.Sprintf("ship_%s_%s", req.OrderID, uuid.New().String()[:8])
	
	// In a real implementation, you'd integrate with shipping providers
	log.Printf("Shipping processed with ID: %s", shippingID)
	
	return shippingID, nil
}