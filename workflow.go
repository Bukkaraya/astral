package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// SimpleWorkflow is a workflow that executes two activities
func SimpleWorkflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result1 string
	err := workflow.ExecuteActivity(ctx, ProcessActivity, name).Get(ctx, &result1)
	if err != nil {
		return "", err
	}

	var result2 string
	err = workflow.ExecuteActivity(ctx, ValidateActivity, result1).Get(ctx, &result2)
	if err != nil {
		return "", err
	}

	return result2, nil
}