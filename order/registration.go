package order

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Orchestrator struct {
	workflows  Workflows
	activities Activities
	client     client.Client
	taskQueue  string
}

func NewOrchestrator(workflows Workflows, activities Activities) *Orchestrator {
	return &Orchestrator{
		workflows:  workflows,
		activities: activities,
	}
}

// SetClient sets the Temporal client and task queue for workflow execution
func (o *Orchestrator) SetClient(client client.Client, taskQueue string) {
	o.client = client
	o.taskQueue = taskQueue
}


func (o *Orchestrator) RegisterWithWorker(w worker.Worker) {
	// Register workflows with versioning
	w.RegisterWorkflowWithOptions(o.workflows.ProcessOrder, workflow.RegisterOptions{
		Name: "ProcessOrder.v1",
	})
	w.RegisterWorkflowWithOptions(o.workflows.CancelOrder, workflow.RegisterOptions{
		Name: "CancelOrder.v1",
	})

	// Register activities
	w.RegisterActivity(o.activities.ValidateOrder)
	w.RegisterActivity(o.activities.ReserveInventory)
	w.RegisterActivity(o.activities.ProcessShipping)
	w.RegisterActivity(o.activities.UpdateOrderStatus)
}