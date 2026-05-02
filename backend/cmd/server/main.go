package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/db"
	"morphic-os/backend/internal/infrastructure/llm"
	"morphic-os/backend/internal/infrastructure/wasm"
	morphichttp "morphic-os/backend/internal/interface/http"
	"morphic-os/backend/internal/usecase"
)

func main() {
	ctx := context.Background()

	// 0. Load Configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}
	cfg, err := domain.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 1. Initialize Database
	toolRepo, err := db.NewSQLiteToolRepository(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	workspaceRepo := db.NewSQLiteWorkspaceRepository(toolRepo.GetDB())

	// 2. Initialize Sandbox
	sandbox := wasm.NewWazeroSandboxManager(ctx)
	defer func() {
		if err := sandbox.Close(ctx); err != nil {
			log.Printf("Failed to close sandbox: %v", err)
		}
	}()

	// 3. Initialize LLM Agent Factory (Abstracted for future expansion)
	agent, err := llm.NewAgentFactory(cfg.Active, cfg.LLM[cfg.Active])
	if err != nil {
		log.Printf("Failed to initialize specified agent '%s': %v. Falling back to MockAgent.", cfg.Active, err)
		agent = llm.NewMockAgent()
	}

	// Initialize Broadcaster
	broadcaster := morphichttp.NewBroadcaster()

	vfsRepo := db.NewSQLiteVirtualFileRepository(toolRepo.GetDB())

	// 4. Initialize UseCase (Morphic Loop)
	morphicLoop := usecase.NewMorphicLoop(toolRepo, workspaceRepo, agent, sandbox)
	morphicLoop.SetLogBroadcaster(broadcaster.Broadcast)

	// 5. Initialize HTTP Handler and Router
	handler := morphichttp.NewHandler(morphicLoop, toolRepo, workspaceRepo, vfsRepo, broadcaster)
	router := morphichttp.SetupRouter(handler)

	// 6. Start HTTP Server
	log.Printf("Starting Morphic-OS server on :%s using %s agent", cfg.Port, cfg.Active)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
