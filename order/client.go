package order

import (
	"context"
	"fmt"

	"simple-temporal-workflow/common"
	"simple-temporal-workflow/order/workflows"
	temporalclient "go.temporal.io/sdk/client"
)

// ProcessOrderRequest represents a request to process an order
type ProcessOrderRequest struct {
	OrderID string `json:"orderId"`
	UserID  string `json:"userId,omitempty"` // Optional for search attributes
}

// CancelOrderRequest represents a request to cancel an order
type CancelOrderRequest struct {
	OrderID string `json:"orderId"`
	UserID  string `json:"userId,omitempty"` // Optional for search attributes
}

// Client provides methods to execute order workflows
type Client interface {
	ProcessOrder(ctx context.Context, req ProcessOrderRequest) (*common.WorkflowResult, error)
	CancelOrder(ctx context.Context, req CancelOrderRequest) (*common.WorkflowResult, error)
}

// orderClient implements the Client interface
type orderClient struct {
	commonClient *common.Client
}

// NewClient creates a new order workflow client
func NewClient(temporalClient temporalclient.Client, taskQueue string) Client {
	return &orderClient{
		commonClient: common.NewClient(temporalClient, taskQueue),
	}
}


// ProcessOrder starts a ProcessOrder workflow
func (c *orderClient) ProcessOrder(ctx context.Context, req ProcessOrderRequest) (*common.WorkflowResult, error) {
	workflowInput := workflows.OrderRequest{OrderID: req.OrderID}
	
	// Build search attributes
	searchAttributes := make(map[string]any)
	if req.UserID != "" {
		searchAttributes["userId"] = req.UserID
	}
	
	return c.commonClient.ExecuteWorkflow(ctx, common.WorkflowExecutionParams{
		WorkflowType:     "ProcessOrder.v1",
		WorkflowIDPrefix: "process-order",
		WorkflowInput:    workflowInput,
		SearchAttributes: searchAttributes,
		SuccessMessage:   fmt.Sprintf("Order processing workflow started for order %s", req.OrderID),
	})
}

// CancelOrder starts a CancelOrder workflow
func (c *orderClient) CancelOrder(ctx context.Context, req CancelOrderRequest) (*common.WorkflowResult, error) {
	workflowInput := workflows.OrderRequest{OrderID: req.OrderID}
	
	// Build search attributes
	searchAttributes := make(map[string]any)
	if req.UserID != "" {
		searchAttributes["userId"] = req.UserID
	}
	
	return c.commonClient.ExecuteWorkflow(ctx, common.WorkflowExecutionParams{
		WorkflowType:     "CancelOrder.v1",
		WorkflowIDPrefix: "cancel-order",
		WorkflowInput:    workflowInput,
		SearchAttributes: searchAttributes,
		SuccessMessage:   fmt.Sprintf("Order cancellation workflow started for order %s", req.OrderID),
	})
}