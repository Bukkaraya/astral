package payment

import (
	"simple-temporal-workflow/payment/workflows"
	"go.temporal.io/sdk/workflow"
)

// paymentWorkflowAdapter adapts the payment workflows to the domain interface
type paymentWorkflowAdapter struct {
	workflows *workflows.Workflows
}

func (a *paymentWorkflowAdapter) ProcessPayment(ctx workflow.Context, req workflows.PaymentRequest) (string, error) {
	return a.workflows.ProcessPayment(ctx, req)
}

func (a *paymentWorkflowAdapter) RefundPayment(ctx workflow.Context, req workflows.RefundRequest) (bool, error) {
	return a.workflows.RefundPayment(ctx, req)
}

// NewWorkflows creates a new payment workflows service
func NewWorkflows(activities Activities) Workflows {
	return &paymentWorkflowAdapter{
		workflows: workflows.NewWorkflows(activities),
	}
}