package workflows

import (
	"errors"
	"testing"

	"simple-temporal-workflow/payment/activities"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

// RefundPaymentTestSuite defines the test suite for RefundPayment workflow
type RefundPaymentTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestRefundPaymentTestSuite(t *testing.T) {
	suite.Run(t, new(RefundPaymentTestSuite))
}

func (s *RefundPaymentTestSuite) TestRefundPayment_Success() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := RefundRequest{
		PaymentID: "payment-123",
	}
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ProcessRefund, mock.Anything, activities.ProcessRefundRequest{PaymentID: req.PaymentID}).Return("refund-456", nil)
	env.OnActivity(mockActivities.UpdatePaymentStatus, mock.Anything, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "refunded"}).Return(nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.RefundPayment, req)
	
	// Verify workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	
	// Verify the result
	var result bool
	s.NoError(env.GetWorkflowResult(&result))
	s.True(result)
}

func (s *RefundPaymentTestSuite) TestRefundPayment_RefundProcessingFailure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := RefundRequest{
		PaymentID: "payment-123",
	}
	refundError := errors.New("refund period expired")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ProcessRefund, mock.Anything, activities.ProcessRefundRequest{PaymentID: req.PaymentID}).Return("", refundError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.RefundPayment, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to process refund")
	s.Contains(env.GetWorkflowError().Error(), "refund period expired")
}

func (s *RefundPaymentTestSuite) TestRefundPayment_StatusUpdateFailure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	req := RefundRequest{
		PaymentID: "payment-123",
	}
	statusError := errors.New("status update service down")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ProcessRefund, mock.Anything, activities.ProcessRefundRequest{PaymentID: req.PaymentID}).Return("refund-456", nil)
	env.OnActivity(mockActivities.UpdatePaymentStatus, mock.Anything, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "refunded"}).Return(statusError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.RefundPayment, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to update payment status")
	s.Contains(env.GetWorkflowError().Error(), "status update service down")
}

func (s *RefundPaymentTestSuite) TestRefundPayment_MultiplePayments() {
	paymentIDs := []string{"payment-1", "payment-2", "payment-3"}
	
	for _, paymentID := range paymentIDs {
		s.Run("refund_"+paymentID, func() {
			env := s.NewTestWorkflowEnvironment()
			
			mockActivities := &MockActivities{}
			workflows := NewWorkflows(mockActivities)
			
			req := RefundRequest{
				PaymentID: paymentID,
			}
			
			// Register mock activities with the test environment
			env.OnActivity(mockActivities.ProcessRefund, mock.Anything, activities.ProcessRefundRequest{PaymentID: req.PaymentID}).Return("refund-"+paymentID, nil)
			env.OnActivity(mockActivities.UpdatePaymentStatus, mock.Anything, activities.UpdatePaymentStatusRequest{PaymentID: req.PaymentID, Status: "refunded"}).Return(nil)
			
			// Execute the workflow
			env.ExecuteWorkflow(workflows.RefundPayment, req)
			
			// Verify workflow completed successfully
			s.True(env.IsWorkflowCompleted())
			s.NoError(env.GetWorkflowError())
			
			var result bool
			s.NoError(env.GetWorkflowResult(&result))
			s.True(result)
		})
	}
}