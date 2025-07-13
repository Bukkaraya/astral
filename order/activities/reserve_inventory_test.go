package activities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivities_ReserveInventory(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	t.Run("successful reservation", func(t *testing.T) {
		start := time.Now()
		req := ReserveInventoryRequest{OrderID: "order-123"}
		reservationID, err := activities.ReserveInventory(ctx, req)
		elapsed := time.Since(start)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, reservationID)
		assert.Contains(t, reservationID, "res_order-123_")
		assert.GreaterOrEqual(t, elapsed, 200*time.Millisecond)
	})

	t.Run("with different order ID", func(t *testing.T) {
		req1 := ReserveInventoryRequest{OrderID: "order-123"}
		req2 := ReserveInventoryRequest{OrderID: "order-456"}
		reservationID1, err1 := activities.ReserveInventory(ctx, req1)
		reservationID2, err2 := activities.ReserveInventory(ctx, req2)
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, reservationID1, reservationID2)
	})
}