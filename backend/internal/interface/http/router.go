package http

import (
	"net/http"
)

// SetupRouter configures the HTTP routes.
func SetupRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/task", handler.HandleTask)
	mux.HandleFunc("/api/tools", handler.HandleGetTools)

	return mux
}
