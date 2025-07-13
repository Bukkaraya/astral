package workflows

import (
	"fmt"

	"simple-temporal-workflow/common"
	"simple-temporal-workflow/order/activities"
	"go.temporal.io/sdk/workflow"
)

func (w *Workflows) CancelOrder(ctx workflow.Context, req OrderRequest) (bool, error) {
	// Apply default activity options
	ctx = common.WithActivityOptions(ctx)

	// Update order status to cancelled
	err := workflow.ExecuteActivity(ctx, w.activities.UpdateOrderStatus, activities.UpdateOrderStatusRequest{OrderID: req.OrderID, Status: "cancelled"}).Get(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to cancel order: %w", err)
	}

	return true, nil
}