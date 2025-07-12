package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessActivityDirect(t *testing.T) {
	ctx := context.Background()
	
	result, err := ProcessActivity(ctx, "test")
	
	assert.NoError(t, err)
	assert.Equal(t, "Processed: test", result)
}

func TestValidateActivityDirect(t *testing.T) {
	ctx := context.Background()
	
	result, err := ValidateActivity(ctx, "input")
	
	assert.NoError(t, err)
	assert.Equal(t, "Validated: input", result)
}