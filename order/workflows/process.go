package workflows

import (
	"context"
	"fmt"

	"simple-temporal-workflow/common"
	"simple-temporal-workflow/order/activities"
	"go.temporal.io/sdk/workflow"
)

// OrderRequest represents an order workflow input
type OrderRequest struct {
	OrderID string `json:"orderId"`
}

// Activities interface for order activities
type Activities interface {
	ValidateOrder(ctx context.Context, req activities.ValidateOrderRequest) (bool, error)
	ReserveInventory(ctx context.Context, req activities.ReserveInventoryRequest) (string, error)
	ProcessShipping(ctx context.Context, req activities.ProcessShippingRequest) (string, error)
	UpdateOrderStatus(ctx context.Context, req activities.UpdateOrderStatusRequest) error
}

func (w *Workflows) ProcessOrder(ctx workflow.Context, req OrderRequest) (string, error) {
	// Apply default activity options
	ctx = common.WithActivityOptions(ctx)

	// Step 1: Validate order
	var isValid bool
	err := workflow.ExecuteActivity(ctx, w.activities.ValidateOrder, activities.ValidateOrderRequest{OrderID: req.OrderID}).Get(ctx, &isValid)
	if err != nil {
		return "", fmt.Errorf("failed to validate order: %w", err)
	}
	if !isValid {
		return "", fmt.Errorf("order validation failed for order %s", req.OrderID)
	}

	// Step 2: Reserve inventory
	var reservationID string
	err = workflow.ExecuteActivity(ctx, w.activities.ReserveInventory, activities.ReserveInventoryRequest{OrderID: req.OrderID}).Get(ctx, &reservationID)
	if err != nil {
		return "", fmt.Errorf("failed to reserve inventory: %w", err)
	}

	// Step 3: Process shipping
	var shippingID string
	err = workflow.ExecuteActivity(ctx, w.activities.ProcessShipping, activities.ProcessShippingRequest{OrderID: req.OrderID}).Get(ctx, &shippingID)
	if err != nil {
		// Compensate: Release inventory reservation
		workflow.ExecuteActivity(ctx, w.activities.UpdateOrderStatus, activities.UpdateOrderStatusRequest{OrderID: req.OrderID, Status: "failed"}).Get(ctx, nil)
		return "", fmt.Errorf("failed to process shipping: %w", err)
	}

	// Step 4: Update order status
	err = workflow.ExecuteActivity(ctx, w.activities.UpdateOrderStatus, activities.UpdateOrderStatusRequest{OrderID: req.OrderID, Status: "completed"}).Get(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	return fmt.Sprintf("Order %s processed successfully. Shipping: %s", req.OrderID, shippingID), nil
}