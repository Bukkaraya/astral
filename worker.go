package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple-temporal-workflow/api"
	"simple-temporal-workflow/order"
	"simple-temporal-workflow/order/activities"
	"simple-temporal-workflow/payment"
	paymentactivities "simple-temporal-workflow/payment/activities"
	myworker "simple-temporal-workflow/worker"
	"syscall"
	"time"

	"go.temporal.io/sdk/worker"
)

const (
	shutdownTimeout = 30 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orderActivities := activities.NewActivities()
	orderWorkflows := order.NewWorkflows(orderActivities)
	orderOrchestrator := order.NewOrchestrator(orderWorkflows, orderActivities)

	paymentActivities := paymentactivities.NewActivities()
	paymentWorkflows := payment.NewWorkflows(paymentActivities)
	paymentOrchestrator := payment.NewOrchestrator(paymentWorkflows, paymentActivities)

	config := myworker.DefaultConfig()
	if taskQueue := os.Getenv("TEMPORAL_TASK_QUEUE"); taskQueue != "" {
		config.TaskQueue = taskQueue
	}
	if hostPort := os.Getenv("TEMPORAL_HOST_PORT"); hostPort != "" {
		config.HostPort = hostPort
	}

	embeddedWorker, err := myworker.NewEmbeddedWorker(config, func(w worker.Worker) {
		orderOrchestrator.RegisterWithWorker(w)
		paymentOrchestrator.RegisterWithWorker(w)
	})
	if err != nil {
		log.Fatalf("Failed to create embedded worker: %v", err)
	}

	if err := embeddedWorker.Start(ctx); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	// Configure orchestrators with client
	orderOrchestrator.SetClient(embeddedWorker.GetClient(), config.TaskQueue)

	// Create domain clients using the same temporal client
	temporalClient := embeddedWorker.GetClient()
	orderClient := order.NewClient(temporalClient, config.TaskQueue)

	// Start API server for workflow triggers
	apiServer := api.NewServer(orderClient)
	mux := http.NewServeMux()
	apiServer.RegisterRoutes(mux)

	apiHttpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Printf("Starting API server on :8080")
		if err := apiHttpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("API server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Microservice started successfully")
	log.Printf("- Task Queue: %s", config.TaskQueue)
	log.Printf("- Temporal Host: %s", config.HostPort)
	log.Printf("- Health endpoint: http://localhost:9090/health")
	log.Printf("- Ready endpoint: http://localhost:9090/ready")
	log.Printf("- Metrics endpoint: http://localhost:9090/metrics")
	log.Printf("- API endpoints: http://localhost:8080/api/workflows/*")

	<-sigChan
	log.Printf("Received shutdown signal, starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Stop API server
	if err := apiHttpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during API server shutdown: %v", err)
	}

	if err := embeddedWorker.Stop(shutdownCtx); err != nil {
		log.Printf("Error during worker shutdown: %v", err)
	}

	log.Printf("Microservice shut down completed")
}