package activities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivities_ChargePayment(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	t.Run("successful charge", func(t *testing.T) {
		start := time.Now()
		transactionID, err := activities.ChargePayment(ctx, ChargePaymentRequest{PaymentID: "payment-123", Amount: 99.99})
		elapsed := time.Since(start)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, transactionID)
		assert.Contains(t, transactionID, "txn_payment-123_")
		assert.GreaterOrEqual(t, elapsed, 500*time.Millisecond)
	})

	t.Run("different amounts", func(t *testing.T) {
		testCases := []struct {
			name      string
			paymentID string
			amount    float64
		}{
			{"small amount", "payment-1", 1.00},
			{"large amount", "payment-2", 9999.99},
			{"decimal amount", "payment-3", 123.45},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				transactionID, err := activities.ChargePayment(ctx, ChargePaymentRequest{PaymentID: tc.paymentID, Amount: tc.amount})
				
				assert.NoError(t, err)
				assert.NotEmpty(t, transactionID)
				assert.Contains(t, transactionID, "txn_"+tc.paymentID+"_")
			})
		}
	})

	t.Run("generates unique transaction IDs", func(t *testing.T) {
		txnID1, err1 := activities.ChargePayment(ctx, ChargePaymentRequest{PaymentID: "payment-123", Amount: 99.99})
		txnID2, err2 := activities.ChargePayment(ctx, ChargePaymentRequest{PaymentID: "payment-123", Amount: 99.99})
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, txnID1, txnID2)
	})
}