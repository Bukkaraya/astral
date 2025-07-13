package payment

import (
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Orchestrator struct {
	workflows  Workflows
	activities Activities
}

func NewOrchestrator(workflows Workflows, activities Activities) *Orchestrator {
	return &Orchestrator{
		workflows:  workflows,
		activities: activities,
	}
}

func (o *Orchestrator) RegisterWithWorker(w worker.Worker) {
	// Register workflows with versioning
	w.RegisterWorkflowWithOptions(o.workflows.ProcessPayment, workflow.RegisterOptions{
		Name: "ProcessPayment.v1",
	})
	w.RegisterWorkflowWithOptions(o.workflows.RefundPayment, workflow.RegisterOptions{
		Name: "RefundPayment.v1",
	})

	// Register activities
	w.RegisterActivity(o.activities.ValidatePayment)
	w.RegisterActivity(o.activities.ChargePayment)
	w.RegisterActivity(o.activities.ProcessRefund)
	w.RegisterActivity(o.activities.UpdatePaymentStatus)
}