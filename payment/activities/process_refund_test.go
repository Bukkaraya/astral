package activities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivities_ProcessRefund(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	t.Run("successful refund", func(t *testing.T) {
		start := time.Now()
		refundID, err := activities.ProcessRefund(ctx, ProcessRefundRequest{PaymentID: "payment-123"})
		elapsed := time.Since(start)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, refundID)
		assert.Contains(t, refundID, "ref_payment-123_")
		assert.GreaterOrEqual(t, elapsed, 400*time.Millisecond)
	})

	t.Run("different payment IDs", func(t *testing.T) {
		paymentIDs := []string{"payment-1", "payment-abc", "payment-xyz-123"}
		
		for _, paymentID := range paymentIDs {
			t.Run("payment_"+paymentID, func(t *testing.T) {
				refundID, err := activities.ProcessRefund(ctx, ProcessRefundRequest{PaymentID: paymentID})
				
				assert.NoError(t, err)
				assert.NotEmpty(t, refundID)
				assert.Contains(t, refundID, "ref_"+paymentID+"_")
			})
		}
	})

	t.Run("generates unique refund IDs", func(t *testing.T) {
		refundID1, err1 := activities.ProcessRefund(ctx, ProcessRefundRequest{PaymentID: "payment-123"})
		refundID2, err2 := activities.ProcessRefund(ctx, ProcessRefundRequest{PaymentID: "payment-123"})
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, refundID1, refundID2)
	})
}