package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"morphic-os/backend/internal/infrastructure/db"
	"morphic-os/backend/internal/infrastructure/llm"
	"morphic-os/backend/internal/infrastructure/wasm"
	morphichttp "morphic-os/backend/internal/interface/http"
	"morphic-os/backend/internal/usecase"
)

func main() {
	ctx := context.Background()

	// 1. Initialize Database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "morphic-os.db"
	}
	toolRepo, err := db.NewSQLiteToolRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. Initialize Sandbox
	sandbox := wasm.NewWazeroSandboxManager(ctx)
	defer func() {
		if err := sandbox.Close(ctx); err != nil {
			log.Printf("Failed to close sandbox: %v", err)
		}
	}()

	// 3. Initialize LLM Agent
	agent := llm.NewMockAgent()

	// Initialize Broadcaster
	broadcaster := morphichttp.NewBroadcaster()
	// Start broadcaster in background (not strictly necessary but good practice if needed)

	// 4. Initialize UseCase (Morphic Loop)
	morphicLoop := usecase.NewMorphicLoop(toolRepo, agent, sandbox)
	morphicLoop.SetLogBroadcaster(broadcaster.Broadcast)

	// 5. Initialize HTTP Handler and Router
	handler := morphichttp.NewHandler(morphicLoop, toolRepo, broadcaster)
	router := morphichttp.SetupRouter(handler)

	// 6. Start HTTP Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Morphic-OS server on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
