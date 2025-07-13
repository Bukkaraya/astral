package activities

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type ReserveInventoryRequest struct {
	OrderID string `json:"orderId"`
}

func (a *Activities) ReserveInventory(ctx context.Context, req ReserveInventoryRequest) (string, error) {
	log.Printf("Reserving inventory for order: %s", req.OrderID)
	
	// Simulate inventory reservation
	time.Sleep(200 * time.Millisecond)
	
	reservationID := fmt.Sprintf("res_%s_%s", req.OrderID, uuid.New().String()[:8])
	
	// In a real implementation, you'd update inventory tables, etc.
	log.Printf("Inventory reserved with ID: %s", reservationID)
	
	return reservationID, nil
}