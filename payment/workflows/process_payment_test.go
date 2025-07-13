package workflows

import (
	"errors"
	"testing"

	"simple-temporal-workflow/payment/activities"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

// ProcessPaymentTestSuite defines the test suite for ProcessPayment workflow
type ProcessPaymentTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestProcessPaymentTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessPaymentTestSuite))
}

func (s *ProcessPaymentTestSuite) TestProcessPayment_Success() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := PaymentRequest{
		PaymentID: "payment-123",
		Amount:    99.99,
	}
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidatePayment, mock.Anything, activities.ValidatePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return(true, nil)
	env.OnActivity(mockActivities.ChargePayment, mock.Anything, activities.ChargePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return("txn-456", nil)
	env.OnActivity(mockActivities.UpdatePaymentStatus, mock.Anything, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "completed"}).Return(nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessPayment, req)
	
	// Verify workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	
	// Verify the result
	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Contains(result, "Payment payment-123 processed successfully")
	s.Contains(result, "Transaction: txn-456")
}

func (s *ProcessPaymentTestSuite) TestProcessPayment_ValidationFailure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := PaymentRequest{
		PaymentID: "invalid-payment",
		Amount:    -10.0,
	}
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidatePayment, mock.Anything, activities.ValidatePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return(false, nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessPayment, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "payment validation failed")
	s.Contains(env.GetWorkflowError().Error(), "invalid-payment")
}

func (s *ProcessPaymentTestSuite) TestProcessPayment_ValidationError() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := PaymentRequest{
		PaymentID: "payment-123",
		Amount:    99.99,
	}
	validationError := errors.New("payment service unavailable")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidatePayment, mock.Anything, activities.ValidatePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return(false, validationError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessPayment, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to validate payment")
	s.Contains(env.GetWorkflowError().Error(), "payment service unavailable")
}

func (s *ProcessPaymentTestSuite) TestProcessPayment_ChargeFailureWithCompensation() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := PaymentRequest{
		PaymentID: "payment-123",
		Amount:    99.99,
	}
	chargeError := errors.New("insufficient funds")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidatePayment, mock.Anything, activities.ValidatePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return(true, nil)
	env.OnActivity(mockActivities.ChargePayment, mock.Anything, activities.ChargePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return("", chargeError)
	// Expect compensation - update status to failed
	env.OnActivity(mockActivities.UpdatePaymentStatus, mock.Anything, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "failed"}).Return(nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessPayment, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to charge payment")
	s.Contains(env.GetWorkflowError().Error(), "insufficient funds")
}

func (s *ProcessPaymentTestSuite) TestProcessPayment_StatusUpdateFailure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := PaymentRequest{
		PaymentID: "payment-123",
		Amount:    99.99,
	}
	statusError := errors.New("database connection lost")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidatePayment, mock.Anything, activities.ValidatePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return(true, nil)
	env.OnActivity(mockActivities.ChargePayment, mock.Anything, activities.ChargePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return("txn-456", nil)
	env.OnActivity(mockActivities.UpdatePaymentStatus, mock.Anything, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "completed"}).Return(statusError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessPayment, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to update payment status")
	s.Contains(env.GetWorkflowError().Error(), "database connection lost")
}

func (s *ProcessPaymentTestSuite) TestProcessPayment_DifferentAmounts() {
	testCases := []struct {
		name      string
		paymentID string
		amount    float64
	}{
		{"small amount", "payment-small", 1.00},
		{"large amount", "payment-large", 9999.99},
		{"decimal amount", "payment-decimal", 123.45},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			env := s.NewTestWorkflowEnvironment()
			
			mockActivities := &MockActivities{}
			workflows := NewWorkflows(mockActivities)
			
			req := PaymentRequest{
				PaymentID: tc.paymentID,
				Amount:    tc.amount,
			}
			
			// Register mock activities with the test environment
			env.OnActivity(mockActivities.ValidatePayment, mock.Anything, activities.ValidatePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return(true, nil)
			env.OnActivity(mockActivities.ChargePayment, mock.Anything, activities.ChargePaymentRequest{PaymentID: req.PaymentID, Amount: req.Amount}).Return("txn-"+tc.paymentID, nil)
			env.OnActivity(mockActivities.UpdatePaymentStatus, mock.Anything, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "completed"}).Return(nil)
			
			// Execute the workflow
			env.ExecuteWorkflow(workflows.ProcessPayment, req)
			
			// Verify workflow completed successfully
			s.True(env.IsWorkflowCompleted())
			s.NoError(env.GetWorkflowError())
			
			var result string
			s.NoError(env.GetWorkflowResult(&result))
			s.Contains(result, tc.paymentID)
		})
	}
}