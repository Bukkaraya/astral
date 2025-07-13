package common

import (
	"context"
	"fmt"
	"time"

	temporalclient "go.temporal.io/sdk/client"
)

// WorkflowExecutionParams holds parameters for executing a workflow
type WorkflowExecutionParams struct {
	WorkflowType      string
	WorkflowIDPrefix  string
	WorkflowInput     any
	SearchAttributes  map[string]any // Flexible search attributes
	SuccessMessage    string
}

// WorkflowResult represents the result of starting a workflow
type WorkflowResult struct {
	WorkflowID string `json:"workflowId"`
	RunID      string `json:"runId"`
	Message    string `json:"message"`
}

// Client provides common workflow execution functionality
type Client struct {
	temporalClient temporalclient.Client
	taskQueue      string
}

// NewClient creates a new common workflow client
func NewClient(temporalClient temporalclient.Client, taskQueue string) *Client {
	return &Client{
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
	}
}

// ExecuteWorkflow starts a workflow with the given parameters
func (c *Client) ExecuteWorkflow(ctx context.Context, params WorkflowExecutionParams) (*WorkflowResult, error) {
	// Generate unique workflow ID  
	workflowID := fmt.Sprintf("%s-%d", params.WorkflowIDPrefix, time.Now().Unix())

	// Setup workflow options
	options := temporalclient.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: c.taskQueue,
	}

	// Add search attributes if provided
	if len(params.SearchAttributes) > 0 {
		options.SearchAttributes = params.SearchAttributes
	}

	// Start workflow
	run, err := c.temporalClient.ExecuteWorkflow(ctx, options, params.WorkflowType, params.WorkflowInput)
	if err != nil {
		return nil, fmt.Errorf("failed to start %s workflow: %w", params.WorkflowType, err)
	}

	return &WorkflowResult{
		WorkflowID: run.GetID(),
		RunID:      run.GetRunID(),
		Message:    params.SuccessMessage,
	}, nil
}