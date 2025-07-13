package workflows

import (
	"context"

	"github.com/stretchr/testify/mock"
	"simple-temporal-workflow/order/activities"
)

// MockActivities implements the Activities interface for testing
type MockActivities struct {
	mock.Mock
}

func (m *MockActivities) ValidateOrder(ctx context.Context, req activities.ValidateOrderRequest) (bool, error) {
	args := m.Called(ctx, req)
	return args.Bool(0), args.Error(1)
}

func (m *MockActivities) ReserveInventory(ctx context.Context, req activities.ReserveInventoryRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockActivities) ProcessShipping(ctx context.Context, req activities.ProcessShippingRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockActivities) UpdateOrderStatus(ctx context.Context, req activities.UpdateOrderStatusRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}