package workflows

import (
	"fmt"

	"simple-temporal-workflow/common"
	"simple-temporal-workflow/payment/activities"
	"go.temporal.io/sdk/workflow"
)

// RefundRequest represents a refund workflow input
type RefundRequest struct {
	PaymentID string `json:"paymentId"`
}

func (w *Workflows) RefundPayment(ctx workflow.Context, req RefundRequest) (bool, error) {
	// Apply default activity options
	ctx = common.WithActivityOptions(ctx)

	// Step 1: Process refund
	var refundID string
	err := workflow.ExecuteActivity(ctx, w.activities.ProcessRefund, activities.ProcessRefundRequest{PaymentID: req.PaymentID}).Get(ctx, &refundID)
	if err != nil {
		return false, fmt.Errorf("failed to process refund: %w", err)
	}

	// Step 2: Update payment status
	err = workflow.ExecuteActivity(ctx, w.activities.UpdatePaymentStatus, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "refunded"}).Get(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to update payment status: %w", err)
	}

	return true, nil
}