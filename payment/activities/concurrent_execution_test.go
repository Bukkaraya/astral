package activities

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivities_ConcurrentExecution(t *testing.T) {
	activities := NewActivities()
	ctx := context.Background()

	t.Run("concurrent payment processing", func(t *testing.T) {
		const numPayments = 5
		results := make(chan string, numPayments)
		errors := make(chan error, numPayments)

		// Process multiple payments concurrently
		for i := 0; i < numPayments; i++ {
			go func(id int) {
				paymentID := fmt.Sprintf("payment-%d", id)
				transactionID, err := activities.ChargePayment(ctx, ChargePaymentRequest{PaymentID: paymentID, Amount: 100.0})
				if err != nil {
					errors <- err
					return
				}
				results <- transactionID
			}(i)
		}

		// Collect results
		var transactionIDs []string
		for i := 0; i < numPayments; i++ {
			select {
			case txnID := <-results:
				transactionIDs = append(transactionIDs, txnID)
			case err := <-errors:
				t.Fatalf("Unexpected error: %v", err)
			}
		}

		// Verify all transaction IDs are unique
		assert.Len(t, transactionIDs, numPayments)
		uniqueIDs := make(map[string]bool)
		for _, id := range transactionIDs {
			assert.False(t, uniqueIDs[id], "Transaction ID %s is not unique", id)
			uniqueIDs[id] = true
		}
	})
}