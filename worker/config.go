package worker

import (
	"time"

	"go.temporal.io/sdk/worker"
)

type Config struct {
	// Temporal connection
	HostPort  string
	Namespace string
	TaskQueue string

	// Worker configuration for low volume
	MaxConcurrentWorkflows int
	MaxConcurrentActivities int

	// Resource limits
	WorkerStopTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		HostPort:  "localhost:7233",
		Namespace: "claude",
		TaskQueue: "microservice-task-queue",

		// Low volume settings - 10 concurrent as specified
		MaxConcurrentWorkflows:  10,
		MaxConcurrentActivities: 10,

		WorkerStopTimeout: 30 * time.Second,
	}
}

func (c *Config) WorkerOptions() worker.Options {
	return worker.Options{
		MaxConcurrentWorkflowTaskPollers:    c.MaxConcurrentWorkflows,
		MaxConcurrentActivityTaskPollers:    c.MaxConcurrentActivities,
		WorkerStopTimeout:                   c.WorkerStopTimeout,

		// Enable worker versioning for safe deployments
		EnableSessionWorker: true,
		
		// Enable structured logging
		EnableLoggingInReplay: false,
	}
}