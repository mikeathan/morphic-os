package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase/prompts"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// SleepCycleConfig defines the thresholds for pruning memories.
type SleepCycleConfig struct {
	MaxLastRecallDays  int // Days since last recall to consider for pruning
	MaxAccessFrequency int // Maximum access frequency to be considered for pruning
}

// NightlySleepCycle represents the background daemon responsible for pruning cognitive memory.
type NightlySleepCycle struct {
	repo          domain.MemoryRepository
	vfsRepo       domain.VirtualFileRepository
	workspaceRepo domain.WorkspaceRepository
	agent         Agent
	config        SleepCycleConfig
	prunedCount   int64 // Atomic counter for pruned vectors
}

// NewNightlySleepCycle creates a new NightlySleepCycle daemon.
func NewNightlySleepCycle(repo domain.MemoryRepository, vfsRepo domain.VirtualFileRepository, workspaceRepo domain.WorkspaceRepository, agent Agent, config SleepCycleConfig) *NightlySleepCycle {
	if config.MaxLastRecallDays == 0 {
		config.MaxLastRecallDays = 30 // default 30 days
	}
	if config.MaxAccessFrequency == 0 {
		config.MaxAccessFrequency = 5 // default 5 accesses
	}
	return &NightlySleepCycle{
		repo:          repo,
		vfsRepo:       vfsRepo,
		workspaceRepo: workspaceRepo,
		agent:         agent,
		config:        config,
	}
}

// Run executes the sleep cycle, finding prunable vectors and deleting them.
func (n *NightlySleepCycle) Run(ctx context.Context) {
	log.Println("[SleepCycle] Starting Nightly Sleep Cycle...")

	// 1. Consolidate raw VFS chat logs into facts/preferences
	if n.workspaceRepo != nil && n.vfsRepo != nil {
		n.consolidateChatLogs(ctx)
	}

	// 2. Prune old vectors
	maxLastRecallTime := time.Now().AddDate(0, 0, -n.config.MaxLastRecallDays)

	vectors, err := n.repo.GetPrunableVectors(ctx, maxLastRecallTime, n.config.MaxAccessFrequency)
	if err != nil {
		log.Printf("[SleepCycle] Error getting prunable vectors: %v\n", err)
		return
	}

	log.Printf("[SleepCycle] Found %d vectors to prune.\n", len(vectors))

	prunedThisRun := 0
	for _, v := range vectors {
		prompt := fmt.Sprintf(prompts.SleepCycleEvaluateMemory, v.Content)

		respStr, err := n.agent.EvaluateTask(ctx, prompt, nil)

		var agentResp AgentResponse
		if err == nil {
			err = json.Unmarshal([]byte(respStr), &agentResp)
		}

		if err != nil {
			log.Printf("[SleepCycle] Agent evaluation failed for vector %s, falling back to normal pruning: %v\n", v.ID, err)
			// Proceed with normal pruning (DISCARD)
		} else if strings.Contains(strings.ToUpper(agentResp.Response), "KEEP") {
			// Update LastRecall to avoid re-evaluating every night immediately
			v.LastRecall = time.Now()
			if err := n.repo.Save(ctx, v); err != nil {
				log.Printf("[SleepCycle] Failed to update kept vector %s: %v\n", v.ID, err)
			} else {
				log.Printf("[SleepCycle] Agent elected to KEEP vector %s.\n", v.ID)
			}
			continue
		}

		// DISCARD
		if err := n.repo.Delete(ctx, v.ID); err != nil {
			log.Printf("[SleepCycle] Failed to delete vector %s: %v\n", v.ID, err)
			continue
		}
		prunedThisRun++
	}

	atomic.AddInt64(&n.prunedCount, int64(prunedThisRun))
	log.Printf("[SleepCycle] Nightly Sleep Cycle completed. Pruned %d vectors.\n", prunedThisRun)
}

func (n *NightlySleepCycle) consolidateChatLogs(ctx context.Context) {
	workspaces, err := n.workspaceRepo.List(ctx)
	if err != nil {
		log.Printf("[SleepCycle] Error listing workspaces for log consolidation: %v\n", err)
		return
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	logFileName := fmt.Sprintf("/var/logs/chat/%s.jsonl", yesterday)

	for _, workspace := range workspaces {
		file, err := n.vfsRepo.GetByPath(ctx, workspace.ID, logFileName)
		if err != nil {
			// Not an error, just no logs for yesterday in this workspace
			continue
		}

		log.Printf("[SleepCycle] Consolidating logs for workspace %s\n", workspace.ID)

		prompt := fmt.Sprintf(prompts.SleepCycleConsolidateLogs, string(file.Content))

		respStr, err := n.agent.EvaluateTask(ctx, prompt, nil)
		if err != nil {
			log.Printf("[SleepCycle] Agent evaluation failed for consolidation in workspace %s: %v\n", workspace.ID, err)
			continue
		}

		var extractedStatements []string

		// If the response is wrapped in AgentResponse structure (due to CloudProviderAgent), try to extract it
		var agentResp AgentResponse
		if unmarshalErr := json.Unmarshal([]byte(respStr), &agentResp); unmarshalErr == nil && agentResp.Action == "direct_response" {
			respStr = agentResp.Response
		}

		if unmarshalErr := json.Unmarshal([]byte(respStr), &extractedStatements); unmarshalErr == nil {
			for _, statement := range extractedStatements {
				newMemory := &domain.MemoryVector{
					ID:              uuid.New().String(),
					WorkspaceID:     workspace.ID,
					Content:         statement,
					AccessFrequency: 1,
					LastRecall:      time.Now(),
					CoreMemory:      false, // Could be determined by agent in future
				}
				_ = n.repo.Save(ctx, newMemory)
			}
			log.Printf("[SleepCycle] Consolidated %d facts/preferences in workspace %s.\n", len(extractedStatements), workspace.ID)
		} else {
			log.Printf("[SleepCycle] Failed to unmarshal extracted statements for workspace %s: %v\n", workspace.ID, unmarshalErr)
		}

		// Purge raw daily logs
		if err := n.vfsRepo.Delete(ctx, file.ID); err != nil {
			log.Printf("[SleepCycle] Failed to delete consolidated log file %s: %v\n", file.ID, err)
		}
	}
}

// GetPrunedCount returns the total number of pruned vectors since startup.
func (n *NightlySleepCycle) GetPrunedCount() int {
	return int(atomic.LoadInt64(&n.prunedCount))
}
