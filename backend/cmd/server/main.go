package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"bufio"
	"strings"

	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/db"
	"morphic-os/backend/internal/infrastructure/llm"
	"morphic-os/backend/internal/infrastructure/wasm"
	morphichttp "morphic-os/backend/internal/interface/http"
	"morphic-os/backend/internal/usecase"
	"time"
)

func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return // Ignore if not found
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
}

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
	agent := llm.NewAgentFactory(cfg)

	// Initialize Broadcaster
	broadcaster := morphichttp.NewBroadcaster()

	vfsRepo := db.NewSQLiteVirtualFileRepository(toolRepo.GetDB())
	secretRepo := db.NewSQLiteSecretRepository(toolRepo.GetDB())
	memoryRepo := db.NewSQLiteMemoryRepository(toolRepo.GetDB())

	// Load env
	loadEnvFile()
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		log.Println("WARNING: ENCRYPTION_KEY is not set. Using fallback key for development.")
		encryptionKey = "this-is-a-super-secret-key-that-must-be-32-bytes"
	}
	secretSvc := usecase.NewSecretService(secretRepo, encryptionKey)

	// 4. Initialize UseCase (Morphic Loop)
	morphicLoop := usecase.NewMorphicLoop(toolRepo, workspaceRepo, agent, sandbox)
	morphicLoop.SetVFSRepo(vfsRepo)
	morphicLoop.SetLogBroadcaster(broadcaster.Broadcast)

	// Initialize Nightly Sleep Cycle Daemon
	sleepCycleConfig := usecase.SleepCycleConfig{
		MaxLastRecallDays:  30,
		MaxAccessFrequency: 5,
	}
	sleepCycleDaemon := usecase.NewNightlySleepCycle(memoryRepo, vfsRepo, workspaceRepo, agent, sleepCycleConfig)

	// Initialize Scheduler
	scheduler := usecase.NewScheduler()
	schedulerCtx, cancelScheduler := context.WithCancel(ctx)

	// Schedule the Sleep Cycle to run every 24 hours
	if err := scheduler.Schedule("nightly-sleep-cycle", 24*time.Hour, sleepCycleDaemon.Run); err != nil {
		log.Printf("Failed to schedule nightly sleep cycle: %v", err)
	}

	scheduler.Start(schedulerCtx)
	defer func() {
		cancelScheduler()
		scheduler.Stop()
	}()

	// 5. Initialize HTTP Handler and Router
	handlerParams := morphichttp.HandlerParams{
		MorphicLoop:   morphicLoop,
		ToolRepo:      toolRepo,
		WorkspaceRepo: workspaceRepo,
		VFSRepo:       vfsRepo,
		SecretSvc:     secretSvc,
		Broadcaster:   broadcaster,
		SleepCycle:    sleepCycleDaemon,
	}
	handler := morphichttp.NewHandler(handlerParams)
	router := morphichttp.SetupRouter(handler)

	// 6. Start HTTP Server
	log.Printf("Starting Morphic-OS server on :%s using %s agent", cfg.Port, cfg.Active)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
