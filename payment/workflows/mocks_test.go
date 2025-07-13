package workflows

import (
	"context"

	"simple-temporal-workflow/payment/activities"
	"github.com/stretchr/testify/mock"
)

// MockActivities implements the Activities interface for testing
type MockActivities struct {
	mock.Mock
}

func (m *MockActivities) ValidatePayment(ctx context.Context, req activities.ValidatePaymentRequest) (bool, error) {
	args := m.Called(ctx, req)
	return args.Bool(0), args.Error(1)
}

func (m *MockActivities) ChargePayment(ctx context.Context, req activities.ChargePaymentRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockActivities) ProcessRefund(ctx context.Context, req activities.ProcessRefundRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockActivities) UpdatePaymentStatus(ctx context.Context, req activities.UpdatePaymentStatusRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}