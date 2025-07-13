package activities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivities_UpdateOrderStatus(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	testCases := []struct {
		name     string
		orderID  string
		status   string
		wantErr  bool
	}{
		{
			name:    "update to completed",
			orderID: "order-123",
			status:  "completed",
			wantErr: false,
		},
		{
			name:    "update to failed",
			orderID: "order-456",
			status:  "failed",
			wantErr: false,
		},
		{
			name:    "update to cancelled",
			orderID: "order-789",
			status:  "cancelled",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			req := UpdateOrderStatusRequest{OrderID: tc.orderID, Status: tc.status}
			err := activities.UpdateOrderStatus(ctx, req)
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