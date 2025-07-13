package api

import (
	"encoding/json"
	"log"
	"net/http"

	"simple-temporal-workflow/order"
)

// Server handles HTTP requests for triggering workflows
type Server struct {
	orderClient order.Client
}

// NewServer creates a new HTTP server for workflow triggers
func NewServer(orderClient order.Client) *Server {
	return &Server{
		orderClient: orderClient,
	}
}

// RegisterRoutes sets up HTTP routes for workflow triggers
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/workflows/order/process", s.handleProcessOrder)
}


func (s *Server) handleProcessOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req order.ProcessOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.OrderID == "" {
		http.Error(w, "orderId is required", http.StatusBadRequest)
		return
	}

	// Execute workflow through client
	result, err := s.orderClient.ProcessOrder(r.Context(), req)
	if err != nil {
		log.Printf("Failed to start ProcessOrder workflow: %v", err)
		http.Error(w, "Failed to start workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}