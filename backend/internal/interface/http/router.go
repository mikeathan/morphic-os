package http

import (
	"net/http"
)

// SetupRouter configures the HTTP routes and applies middleware.
func SetupRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/task", handler.HandleTask)
	mux.HandleFunc("/api/tools", handler.HandleGetTools)
	mux.HandleFunc("/api/logs", handler.HandleLogs)
	mux.HandleFunc("/api/workspaces", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.HandleGetWorkspaces(w, r)
		} else if r.Method == http.MethodPost {
			handler.HandleCreateWorkspace(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/vfs/files", handler.HandleListFiles)
	mux.HandleFunc("/api/vfs/files/", handler.HandleGetFile)
	mux.HandleFunc("/api/metrics", handler.HandleGetMetrics)

	// Wrap the mux with the CORS middleware
	return CORSMiddleware(mux)
}
