package order

import (
	"context"

	"simple-temporal-workflow/order/activities"
	"simple-temporal-workflow/order/workflows"
	"go.temporal.io/sdk/workflow"
)

// Workflows defines the order workflow interface
type Workflows interface {
	ProcessOrder(ctx workflow.Context, req workflows.OrderRequest) (string, error)
	CancelOrder(ctx workflow.Context, req workflows.OrderRequest) (bool, error)
}

// Activities defines the order activity interface
type Activities interface {
	ValidateOrder(ctx context.Context, req activities.ValidateOrderRequest) (bool, error)
	ReserveInventory(ctx context.Context, req activities.ReserveInventoryRequest) (string, error)
	ProcessShipping(ctx context.Context, req activities.ProcessShippingRequest) (string, error)
	UpdateOrderStatus(ctx context.Context, req activities.UpdateOrderStatusRequest) error
}