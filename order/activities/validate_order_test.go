package activities

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivities_ValidateOrder(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	t.Run("valid order", func(t *testing.T) {
		req := ValidateOrderRequest{OrderID: "order-123"}
		result, err := activities.ValidateOrder(ctx, req)
		
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("empty order ID", func(t *testing.T) {
		req := ValidateOrderRequest{OrderID: ""}
		result, err := activities.ValidateOrder(ctx, req)
		
		assert.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), "order ID cannot be empty")
	})
}