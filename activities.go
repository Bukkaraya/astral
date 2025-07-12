package main

import (
	"context"
	"fmt"
)

// ProcessActivity is the first activity that processes the input
func ProcessActivity(ctx context.Context, name string) (string, error) {
	processed := fmt.Sprintf("Processed: %s", name)
	return processed, nil
}

// ValidateActivity is the second activity that validates the processed input
func ValidateActivity(ctx context.Context, input string) (string, error) {
	validated := fmt.Sprintf("Validated: %s", input)
	return validated, nil
}