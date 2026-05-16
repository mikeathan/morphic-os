package http

import (
	"encoding/json"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
	"net/http"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

// Handler represents the HTTP handler.
// HandlerParams encapsulates all dependencies required by the Handler.
type HandlerParams struct {
	MorphicLoop   *usecase.MorphicLoop
	ToolRepo      domain.ToolRepository
	WorkspaceRepo domain.WorkspaceRepository
	VFSRepo       domain.VirtualFileRepository
	SecretSvc     *usecase.SecretService
	Broadcaster   *Broadcaster
	SleepCycle    *usecase.NightlySleepCycle
	MemoryRepo    domain.MemoryRepository
}

type Handler struct {
	morphicLoop   *usecase.MorphicLoop
	toolRepo      domain.ToolRepository
	workspaceRepo domain.WorkspaceRepository
	vfsRepo       domain.VirtualFileRepository
	secretSvc     *usecase.SecretService
	broadcaster   *Broadcaster
	sleepCycle    *usecase.NightlySleepCycle
	memoryRepo    domain.MemoryRepository
}

// NewHandler creates a new Handler.
func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		morphicLoop:   params.MorphicLoop,
		toolRepo:      params.ToolRepo,
		workspaceRepo: params.WorkspaceRepo,
		vfsRepo:       params.VFSRepo,
		secretSvc:     params.SecretSvc,
		broadcaster:   params.Broadcaster,
		sleepCycle:    params.SleepCycle,
		memoryRepo:    params.MemoryRepo,
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

// HandleListFiles returns a list of virtual files for a workspace.
func (h *Handler) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspaceID := r.URL.Query().Get("workspace_id")
	if workspaceID == "" {
		workspaceID = "default"
	}

	if h.vfsRepo == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]domain.VirtualFile{})
		return
	}

	files, err := h.vfsRepo.ListByWorkspace(r.Context(), workspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Exclude content for listing to save bandwidth
	for i := range files {
		files[i].Content = nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// HandleGetFile returns a single virtual file including its content.
func (h *Handler) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple path parsing since we are not using an advanced router yet
	id := r.URL.Path[len("/api/vfs/files/"):]
	if id == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	if h.vfsRepo == nil {
		http.Error(w, "VFS not configured", http.StatusInternalServerError)
		return
	}

	file, err := h.vfsRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

// HandleGetMetrics returns system metrics.
func (h *Handler) HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := struct {
		Goroutines   int    `json:"goroutines"`
		AllocatedMem uint64 `json:"allocated_mem"`
		TotalAlloc   uint64 `json:"total_alloc"`
		SysMem       uint64 `json:"sys_mem"`
		NumGC        uint32 `json:"num_gc"`
		PrunedCount  int    `json:"pruned_count"` // Placeholder for pruning stats
	}{
		Goroutines:   runtime.NumGoroutine(),
		AllocatedMem: m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		SysMem:       m.Sys,
		NumGC:        m.NumGC,
		PrunedCount:  0,
	}

	if h.sleepCycle != nil {
		metrics.PrunedCount = h.sleepCycle.GetPrunedCount()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.URL.Query().Get("workspace_id")
	if workspaceID == "" {
		workspaceID = "default"
	}

	secrets, err := h.secretSvc.ListSecrets(r.Context(), workspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secrets)
}

func (h *Handler) AddSecret(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WorkspaceID string `json:"workspace_id"`
		Key         string `json:"key"`
		Value       string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.WorkspaceID == "" {
		req.WorkspaceID = "default"
	}

	secret, err := h.secretSvc.AddSecret(r.Context(), req.WorkspaceID, req.Key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secret)
}

func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/secrets/")
	if id == "" {
		http.Error(w, "Missing secret ID", http.StatusBadRequest)
		return
	}

	if err := h.secretSvc.DeleteSecret(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetMemory returns all memory vectors for a workspace.
func (h *Handler) HandleGetMemory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspaceID := r.URL.Query().Get("workspace_id")
	if workspaceID == "" {
		workspaceID = "default"
	}

	if h.memoryRepo == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]domain.MemoryVector{})
		return
	}

	memories, err := h.memoryRepo.ListAll(r.Context(), workspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memories)
}
