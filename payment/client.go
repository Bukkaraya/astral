package payment

import (
	"context"
	"fmt"

	"simple-temporal-workflow/common"
	"simple-temporal-workflow/payment/workflows"
	temporalclient "go.temporal.io/sdk/client"
)

// ProcessPaymentRequest represents a request to process a payment
type ProcessPaymentRequest struct {
	PaymentID string  `json:"paymentId"`
	Amount    float64 `json:"amount"`
	UserID    string  `json:"userId,omitempty"` // Optional for search attributes
}

// RefundPaymentRequest represents a request to refund a payment
type RefundPaymentRequest struct {
	PaymentID string `json:"paymentId"`
	UserID    string `json:"userId,omitempty"` // Optional for search attributes
}

// Client provides methods to execute payment workflows
type Client interface {
	ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*common.WorkflowResult, error)
	RefundPayment(ctx context.Context, req RefundPaymentRequest) (*common.WorkflowResult, error)
}

// paymentClient implements the Client interface
type paymentClient struct {
	commonClient *common.Client
}

// NewClient creates a new payment workflow client
func NewClient(temporalClient temporalclient.Client, taskQueue string) Client {
	return &paymentClient{
		commonClient: common.NewClient(temporalClient, taskQueue),
	}
}

// ProcessPayment starts a ProcessPayment workflow
func (c *paymentClient) ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*common.WorkflowResult, error) {
	workflowInput := workflows.PaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}
	
	// Build search attributes
	searchAttributes := make(map[string]any)
	if req.UserID != "" {
		searchAttributes["userId"] = req.UserID
	}
	
	return c.commonClient.ExecuteWorkflow(ctx, common.WorkflowExecutionParams{
		WorkflowType:     "ProcessPayment.v1",
		WorkflowIDPrefix: "process-payment",
		WorkflowInput:    workflowInput,
		SearchAttributes: searchAttributes,
		SuccessMessage:   fmt.Sprintf("Payment processing workflow started for payment %s", req.PaymentID),
	})
}

// RefundPayment starts a RefundPayment workflow
func (c *paymentClient) RefundPayment(ctx context.Context, req RefundPaymentRequest) (*common.WorkflowResult, error) {
	workflowInput := workflows.RefundRequest{PaymentID: req.PaymentID}
	
	// Build search attributes
	searchAttributes := make(map[string]any)
	if req.UserID != "" {
		searchAttributes["userId"] = req.UserID
	}
	
	return c.commonClient.ExecuteWorkflow(ctx, common.WorkflowExecutionParams{
		WorkflowType:     "RefundPayment.v1",
		WorkflowIDPrefix: "refund-payment",
		WorkflowInput:    workflowInput,
		SearchAttributes: searchAttributes,
		SuccessMessage:   fmt.Sprintf("Payment refund workflow started for payment %s", req.PaymentID),
	})
}