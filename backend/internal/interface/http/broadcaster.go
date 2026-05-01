package http

import (
	"fmt"
	"net/http"
	"sync"
)

// Broadcaster manages SSE clients and broadcasts messages to them.
type Broadcaster struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

// NewBroadcaster creates a new Broadcaster instance.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan string]bool),
	}
}

// AddClient registers a new client channel.
func (b *Broadcaster) AddClient(clientChan chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[clientChan] = true
}

// RemoveClient unregisters a client channel.
func (b *Broadcaster) RemoveClient(clientChan chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[clientChan]; ok {
		delete(b.clients, clientChan)
		close(clientChan)
	}
}

// Broadcast sends a message to all connected clients.
func (b *Broadcaster) Broadcast(message string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for clientChan := range b.clients {
		// Non-blocking send to avoid hanging on slow clients
		select {
		case clientChan <- message:
		default:
			// If the channel is full, we might want to log or disconnect the client
			// For simplicity, we drop the message if the client is too slow.
			fmt.Printf("Warning: Dropping message to slow client\n")
		}
	}
}

// HandleSSE is an HTTP handler for Server-Sent Events.
func (b *Broadcaster) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a channel for this client
	clientChan := make(chan string, 100)
	b.AddClient(clientChan)

	// Clean up on disconnect
	defer b.RemoveClient(clientChan)

	// Notify the client that connection is established
	fmt.Fprintf(w, "data: {\"type\": \"connected\", \"message\": \"SSE connected\"}\n\n")
	flusher.Flush()

	// Listen for client disconnect
	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return
		case msg := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
