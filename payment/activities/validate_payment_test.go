package activities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivities_ValidatePayment(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	t.Run("valid payment", func(t *testing.T) {
		result, err := activities.ValidatePayment(ctx, ValidatePaymentRequest{PaymentID: "payment-123", Amount: 99.99})
		
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("empty payment ID", func(t *testing.T) {
		result, err := activities.ValidatePayment(ctx, ValidatePaymentRequest{PaymentID: "", Amount: 99.99})
		
		assert.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), "payment ID cannot be empty")
	})

	t.Run("zero amount", func(t *testing.T) {
		result, err := activities.ValidatePayment(ctx, ValidatePaymentRequest{PaymentID: "payment-123", Amount: 0})
		
		assert.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), "amount must be positive")
	})

	t.Run("negative amount", func(t *testing.T) {
		result, err := activities.ValidatePayment(ctx, ValidatePaymentRequest{PaymentID: "payment-123", Amount: -10.50})
		
		assert.Error(t, err)
		assert.False(t, result)
		assert.Contains(t, err.Error(), "amount must be positive")
	})

	t.Run("processing time", func(t *testing.T) {
		start := time.Now()
		_, err := activities.ValidatePayment(ctx, ValidatePaymentRequest{PaymentID: "payment-123", Amount: 99.99})
		elapsed := time.Since(start)
		
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond)
	})
}