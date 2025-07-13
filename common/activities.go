package common

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// DefaultActivityOptions returns common activity options used across workflows
func DefaultActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
}

// WithActivityOptions applies default activity options to a workflow context
func WithActivityOptions(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, DefaultActivityOptions())
}