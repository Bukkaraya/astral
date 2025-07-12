package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both workflow and activity functions
	w := worker.New(c, "simple-workflow-task-queue", worker.Options{})

	// Register workflow and activities
	w.RegisterWorkflow(SimpleWorkflow)
	w.RegisterActivity(ProcessActivity)
	w.RegisterActivity(ValidateActivity)

	// Start listening to the task queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}