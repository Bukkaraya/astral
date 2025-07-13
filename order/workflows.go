package order

import (
	"simple-temporal-workflow/order/workflows"
	"go.temporal.io/sdk/workflow"
)

// orderWorkflowAdapter adapts the order workflows to the domain interface
type orderWorkflowAdapter struct {
	workflows *workflows.Workflows
}

func (a *orderWorkflowAdapter) ProcessOrder(ctx workflow.Context, req workflows.OrderRequest) (string, error) {
	return a.workflows.ProcessOrder(ctx, req)
}

func (a *orderWorkflowAdapter) CancelOrder(ctx workflow.Context, req workflows.OrderRequest) (bool, error) {
	return a.workflows.CancelOrder(ctx, req)
}

// NewWorkflows creates a new order workflows service
func NewWorkflows(activities Activities) Workflows {
	return &orderWorkflowAdapter{
		workflows: workflows.NewWorkflows(activities),
	}
}