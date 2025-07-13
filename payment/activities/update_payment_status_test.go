package activities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivities_UpdatePaymentStatus(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	testCases := []struct {
		name      string
		paymentID string
		status    string
		wantErr   bool
	}{
		{
			name:      "update to completed",
			paymentID: "payment-123",
			status:    "completed",
			wantErr:   false,
		},
		{
			name:      "update to failed",
			paymentID: "payment-456",
			status:    "failed",
			wantErr:   false,
		},
		{
			name:      "update to refunded",
			paymentID: "payment-789",
			status:    "refunded",
			wantErr:   false,
		},
		{
			name:      "update to pending",
			paymentID: "payment-abc",
			status:    "pending",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			err := activities.UpdatePaymentStatus(ctx, UpdatePaymentStatusRequest{PaymentID: tc.paymentID, Status: tc.status})
			elapsed := time.Since(start)
			
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
			}
		})
	}
}