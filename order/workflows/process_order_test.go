package workflows

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"simple-temporal-workflow/order/activities"
)

// ProcessOrderTestSuite defines the test suite for ProcessOrder workflow
type ProcessOrderTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestProcessOrderTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessOrderTestSuite))
}

func (s *ProcessOrderTestSuite) TestProcessOrder_Success() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "test-order-123"
	req := OrderRequest{OrderID: orderID}
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidateOrder, mock.Anything, activities.ValidateOrderRequest{OrderID: orderID}).Return(true, nil)
	env.OnActivity(mockActivities.ReserveInventory, mock.Anything, activities.ReserveInventoryRequest{OrderID: orderID}).Return("reservation-123", nil)
	env.OnActivity(mockActivities.ProcessShipping, mock.Anything, activities.ProcessShippingRequest{OrderID: orderID}).Return("shipping-456", nil)
	env.OnActivity(mockActivities.UpdateOrderStatus, mock.Anything, activities.UpdateOrderStatusRequest{OrderID: orderID, Status: "completed"}).Return(nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessOrder, req)
	
	// Verify workflow completed successfully
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	
	// Verify the result
	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Contains(result, "Order test-order-123 processed successfully")
	s.Contains(result, "Shipping: shipping-456")
}

func (s *ProcessOrderTestSuite) TestProcessOrder_ValidationFailure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "invalid-order"
	req := OrderRequest{OrderID: orderID}
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidateOrder, mock.Anything, activities.ValidateOrderRequest{OrderID: orderID}).Return(false, nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessOrder, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "order validation failed")
}

func (s *ProcessOrderTestSuite) TestProcessOrder_ValidationError() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "error-order"
	req := OrderRequest{OrderID: orderID}
	validationError := errors.New("database connection failed")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidateOrder, mock.Anything, activities.ValidateOrderRequest{OrderID: orderID}).Return(false, validationError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessOrder, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to validate order")
	s.Contains(env.GetWorkflowError().Error(), "database connection failed")
}

func (s *ProcessOrderTestSuite) TestProcessOrder_InventoryReservationFailure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "test-order-123"
	req := OrderRequest{OrderID: orderID}
	inventoryError := errors.New("insufficient inventory")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidateOrder, mock.Anything, activities.ValidateOrderRequest{OrderID: orderID}).Return(true, nil)
	env.OnActivity(mockActivities.ReserveInventory, mock.Anything, activities.ReserveInventoryRequest{OrderID: orderID}).Return("", inventoryError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessOrder, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to reserve inventory")
	s.Contains(env.GetWorkflowError().Error(), "insufficient inventory")
}

func (s *ProcessOrderTestSuite) TestProcessOrder_ShippingFailureWithCompensation() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "test-order-123"
	req := OrderRequest{OrderID: orderID}
	shippingError := errors.New("shipping provider unavailable")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidateOrder, mock.Anything, activities.ValidateOrderRequest{OrderID: orderID}).Return(true, nil)
	env.OnActivity(mockActivities.ReserveInventory, mock.Anything, activities.ReserveInventoryRequest{OrderID: orderID}).Return("reservation-123", nil)
	env.OnActivity(mockActivities.ProcessShipping, mock.Anything, activities.ProcessShippingRequest{OrderID: orderID}).Return("", shippingError)
	// Expect compensation - update status to failed
	env.OnActivity(mockActivities.UpdateOrderStatus, mock.Anything, activities.UpdateOrderStatusRequest{OrderID: orderID, Status: "failed"}).Return(nil)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessOrder, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to process shipping")
	s.Contains(env.GetWorkflowError().Error(), "shipping provider unavailable")
}

func (s *ProcessOrderTestSuite) TestProcessOrder_StatusUpdateFailure() {
	env := s.NewTestWorkflowEnvironment()
	
	mockActivities := &MockActivities{}
	workflows := NewWorkflows(mockActivities)
	
	orderID := "test-order-123"
	req := OrderRequest{OrderID: orderID}
	statusError := errors.New("database update failed")
	
	// Register mock activities with the test environment
	env.OnActivity(mockActivities.ValidateOrder, mock.Anything, activities.ValidateOrderRequest{OrderID: orderID}).Return(true, nil)
	env.OnActivity(mockActivities.ReserveInventory, mock.Anything, activities.ReserveInventoryRequest{OrderID: orderID}).Return("reservation-123", nil)
	env.OnActivity(mockActivities.ProcessShipping, mock.Anything, activities.ProcessShippingRequest{OrderID: orderID}).Return("shipping-456", nil)
	env.OnActivity(mockActivities.UpdateOrderStatus, mock.Anything, activities.UpdateOrderStatusRequest{OrderID: orderID, Status: "completed"}).Return(statusError)
	
	// Execute the workflow
	env.ExecuteWorkflow(workflows.ProcessOrder, req)
	
	// Verify workflow completed with error
	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Contains(env.GetWorkflowError().Error(), "failed to update order status")
	s.Contains(env.GetWorkflowError().Error(), "database update failed")
}