package activities

import (
	"context"
	"log"
	"time"
)

type UpdateOrderStatusRequest struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}

func (a *Activities) UpdateOrderStatus(ctx context.Context, req UpdateOrderStatusRequest) error {
	log.Printf("Updating order %s status to: %s", req.OrderID, req.Status)
	
	// Simulate database update
	time.Sleep(50 * time.Millisecond)
	
	// In a real implementation, you'd update your order database
	log.Printf("Order %s status updated to: %s", req.OrderID, req.Status)
	
	return nil
}