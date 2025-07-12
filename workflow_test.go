package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func (s *UnitTestSuite) TestSimpleWorkflow() {
	env := s.NewTestWorkflowEnvironment()
	
	env.OnActivity(ProcessActivity, mock.Anything, "test").Return("Processed: test", nil)
	env.OnActivity(ValidateActivity, mock.Anything, "Processed: test").Return("Validated: Processed: test", nil)
	
	env.ExecuteWorkflow(SimpleWorkflow, "test")
	
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	
	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Equal("Validated: Processed: test", result)
}

func (s *UnitTestSuite) TestSimpleWorkflowWithEmptyName() {
	env := s.NewTestWorkflowEnvironment()
	
	env.OnActivity(ProcessActivity, mock.Anything, "").Return("Processed: ", nil)
	env.OnActivity(ValidateActivity, mock.Anything, "Processed: ").Return("Validated: Processed: ", nil)
	
	env.ExecuteWorkflow(SimpleWorkflow, "")
	
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
	
	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Equal("Validated: Processed: ", result)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}