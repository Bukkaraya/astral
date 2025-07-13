package workflows

import (
	"context"
	"fmt"

	"simple-temporal-workflow/common"
	"simple-temporal-workflow/payment/activities"
	"go.temporal.io/sdk/workflow"
)

// PaymentRequest represents a payment workflow input
type PaymentRequest struct {
	PaymentID string  `json:"paymentId"`
	Amount    float64 `json:"amount"`
}

// Activities interface for payment activities
type Activities interface {
	ValidatePayment(ctx context.Context, req activities.ValidatePaymentRequest) (bool, error)
	ChargePayment(ctx context.Context, req activities.ChargePaymentRequest) (string, error)
	ProcessRefund(ctx context.Context, req activities.ProcessRefundRequest) (string, error)
	UpdatePaymentStatus(ctx context.Context, req activities.UpdatePaymentStatusRequest) error
}

func (w *Workflows) ProcessPayment(ctx workflow.Context, req PaymentRequest) (string, error) {
	// Apply default activity options
	ctx = common.WithActivityOptions(ctx)

	// Step 1: Validate payment
	var isValid bool
	err := workflow.ExecuteActivity(ctx, w.activities.ValidatePayment, activities.ValidatePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Get(ctx, &isValid)
	if err != nil {
		return "", fmt.Errorf("failed to validate payment: %w", err)
	}
	if !isValid {
		return "", fmt.Errorf("payment validation failed for payment %s", req.PaymentID)
	}

	// Step 2: Charge payment
	var transactionID string
	err = workflow.ExecuteActivity(ctx, w.activities.ChargePayment, activities.ChargePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Get(ctx, &transactionID)
	if err != nil {
		// Update payment status to failed
		workflow.ExecuteActivity(ctx, w.activities.UpdatePaymentStatus, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "failed"}).Get(ctx, nil)
		return "", fmt.Errorf("failed to charge payment: %w", err)
	}

	// Step 3: Update payment status
	err = workflow.ExecuteActivity(ctx, w.activities.UpdatePaymentStatus, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "completed"}).Get(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to update payment status: %w", err)
	}

	return fmt.Sprintf("Payment %s processed successfully. Transaction: %s", req.PaymentID, transactionID), nil
}