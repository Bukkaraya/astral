package activities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivities_ProcessShipping(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	t.Run("successful shipping", func(t *testing.T) {
		start := time.Now()
		req := ProcessShippingRequest{OrderID: "order-123"}
		shippingID, err := activities.ProcessShipping(ctx, req)
		elapsed := time.Since(start)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, shippingID)
		assert.Contains(t, shippingID, "ship_order-123_")
		assert.GreaterOrEqual(t, elapsed, 300*time.Millisecond)
	})

	t.Run("generates unique shipping IDs", func(t *testing.T) {
		req := ProcessShippingRequest{OrderID: "order-123"}
		shippingID1, err1 := activities.ProcessShipping(ctx, req)
		shippingID2, err2 := activities.ProcessShipping(ctx, req)
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, shippingID1, shippingID2)
	})
}