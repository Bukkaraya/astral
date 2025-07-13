package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type RegistrationFunc func(worker.Worker)

type EmbeddedWorker struct {
	config           *Config
	client           client.Client
	worker           worker.Worker
	registrationFunc RegistrationFunc
	httpServer       *http.Server
	mu               sync.RWMutex
	running          bool
}

func NewEmbeddedWorker(config *Config, registrationFunc RegistrationFunc) (*EmbeddedWorker, error) {
	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort:  config.HostPort,
		Namespace: config.Namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	// Create worker
	w := worker.New(c, config.TaskQueue, config.WorkerOptions())

	// Register workflows and activities using the provided function
	registrationFunc(w)

	return &EmbeddedWorker{
		config:           config,
		client:           c,
		worker:           w,
		registrationFunc: registrationFunc,
	}, nil
}

func (w *EmbeddedWorker) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return fmt.Errorf("worker is already running")
	}

	// Start health check server
	w.startHealthServer()

	// Start the worker
	log.Printf("Starting Temporal worker on task queue: %s", w.config.TaskQueue)
	
	go func() {
		err := w.worker.Run(worker.InterruptCh())
		if err != nil {
			log.Printf("Worker error: %v", err)
		}
	}()

	w.running = true
	log.Printf("Temporal worker started successfully")
	
	return nil
}

func (w *EmbeddedWorker) Stop(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return nil
	}

	log.Printf("Stopping Temporal worker...")

	// Stop worker gracefully
	w.worker.Stop()

	// Stop health server
	if w.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		w.httpServer.Shutdown(shutdownCtx)
	}

	// Close Temporal client
	w.client.Close()

	w.running = false
	log.Printf("Temporal worker stopped")
	
	return nil
}

func (w *EmbeddedWorker) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.running
}

func (w *EmbeddedWorker) GetClient() client.Client {
	return w.client
}

func (w *EmbeddedWorker) startHealthServer() {
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		if w.IsRunning() {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`{"status":"healthy","worker":"running"}`))
		} else {
			rw.WriteHeader(http.StatusServiceUnavailable)
			rw.Write([]byte(`{"status":"unhealthy","worker":"stopped"}`))
		}
	})

	// Ready check endpoint
	mux.HandleFunc("/ready", func(rw http.ResponseWriter, r *http.Request) {
		if w.IsRunning() {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`{"status":"ready"}`))
		} else {
			rw.WriteHeader(http.StatusServiceUnavailable)
			rw.Write([]byte(`{"status":"not_ready"}`))
		}
	})

	// Metrics endpoint (basic)
	mux.HandleFunc("/metrics", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("temporal_worker_running %d\n", map[bool]int{true: 1, false: 0}[w.IsRunning()])))
	})

	w.httpServer = &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	go func() {
		log.Printf("Starting health server on :9090")
		if err := w.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Health server error: %v", err)
		}
	}()
}