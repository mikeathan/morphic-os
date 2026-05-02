package http

import (
	"encoding/json"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
	"net/http"

	"github.com/google/uuid"
)

// Handler represents the HTTP handler.
type Handler struct {
	morphicLoop   *usecase.MorphicLoop
	toolRepo      domain.ToolRepository
	workspaceRepo domain.WorkspaceRepository
	broadcaster   *Broadcaster
}

// NewHandler creates a new Handler.
func NewHandler(morphicLoop *usecase.MorphicLoop, toolRepo domain.ToolRepository, workspaceRepo domain.WorkspaceRepository, broadcaster *Broadcaster) *Handler {
	return &Handler{
		morphicLoop:   morphicLoop,
		toolRepo:      toolRepo,
		workspaceRepo: workspaceRepo,
		broadcaster:   broadcaster,
	}
}

// HandleLogs handles Server-Sent Events for real-time logs.
func (h *Handler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.broadcaster.HandleSSE(w, r)
}

// HandleTask processes a user task.
func (h *Handler) HandleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Task        string `json:"task"`
		WorkspaceID string `json:"workspace_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Task == "" {
		http.Error(w, "Task is required", http.StatusBadRequest)
		return
	}

	if req.WorkspaceID == "" {
		req.WorkspaceID = "default" // default workspace for now if not provided
	}

	result, err := h.morphicLoop.Execute(r.Context(), req.WorkspaceID, req.Task)
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

	workspaceID := r.URL.Query().Get("workspace_id")
	if workspaceID == "" {
		workspaceID = "default"
	}

	tools, err := h.toolRepo.ListActive(r.Context(), workspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tools)
}

// HandleCreateWorkspace creates a new workspace.
func (h *Handler) HandleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.Workspace
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Workspace Name is required", http.StatusBadRequest)
		return
	}

	// Generate ID server-side to prevent path traversal
	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	if h.workspaceRepo != nil {
		if err := h.workspaceRepo.Create(r.Context(), &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// HandleGetWorkspaces returns a list of workspaces.
func (h *Handler) HandleGetWorkspaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.workspaceRepo == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]domain.Workspace{})
		return
	}

	workspaces, err := h.workspaceRepo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaces)
}
