package payment

import (
	"context"

	"simple-temporal-workflow/payment/activities"
	"simple-temporal-workflow/payment/workflows"
	"go.temporal.io/sdk/workflow"
)

// Workflows defines the payment workflow interface
type Workflows interface {
	ProcessPayment(ctx workflow.Context, req workflows.PaymentRequest) (string, error)
	RefundPayment(ctx workflow.Context, req workflows.RefundRequest) (bool, error)
}

// Activities defines the payment activity interface
type Activities interface {
	ValidatePayment(ctx context.Context, req activities.ValidatePaymentRequest) (bool, error)
	ChargePayment(ctx context.Context, req activities.ChargePaymentRequest) (string, error)
	ProcessRefund(ctx context.Context, req activities.ProcessRefundRequest) (string, error)
	UpdatePaymentStatus(ctx context.Context, req activities.UpdatePaymentStatusRequest) error
}