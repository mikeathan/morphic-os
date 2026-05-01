package http

import (
	"encoding/json"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
	"net/http"
)

// Handler represents the HTTP handler.
type Handler struct {
	morphicLoop *usecase.MorphicLoop
	toolRepo    domain.ToolRepository
}

// NewHandler creates a new Handler.
func NewHandler(morphicLoop *usecase.MorphicLoop, toolRepo domain.ToolRepository) *Handler {
	return &Handler{
		morphicLoop: morphicLoop,
		toolRepo:    toolRepo,
	}
}

// HandleTask processes a user task.
func (h *Handler) HandleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Task string `json:"task"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Task == "" {
		http.Error(w, "Task is required", http.StatusBadRequest)
		return
	}

	result, err := h.morphicLoop.Execute(r.Context(), req.Task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Result string `json:"result"`
	}{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleGetTools returns a list of active tools.
func (h *Handler) HandleGetTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tools, err := h.toolRepo.ListActive(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tools)
}
