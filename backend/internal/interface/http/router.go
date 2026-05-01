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

	// Wrap the mux with the CORS middleware
	return CORSMiddleware(mux)
}
