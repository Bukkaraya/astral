package workflows

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"simple-temporal-workflow/order/activities"
)

// CancelOrderTestSuite defines the test suite for CancelOrder workflow
type CancelOrderTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestCancelOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CancelOrderTestSuite))
}

func (s *CancelOrderTestSuite) TestCancelOrder_Success() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "test-order-123"
	req := OrderRequest{OrderID: orderID}
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.UpdateOrderStatus, mock.Anything, activities.UpdateOrderStatusRequest{OrderID: orderID, Status: "cancelled"}).Return(nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.CancelOrder, req)
	
	// Verify workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	
	// Verify the result
	var result bool
	s.NoError(env.GetWorkflowResult(&result))
	s.True(result)
}

func (s *CancelOrderTestSuite) TestCancelOrder_Failure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "test-order-123"
	req := OrderRequest{OrderID: orderID}
	cancelError := errors.New("order cannot be cancelled")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.UpdateOrderStatus, mock.Anything, activities.UpdateOrderStatusRequest{OrderID: orderID, Status: "cancelled"}).Return(cancelError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.CancelOrder, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to cancel order")
	s.Contains(env.GetWorkflowError().Error(), "order cannot be cancelled")
}