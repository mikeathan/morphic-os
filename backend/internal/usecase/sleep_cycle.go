package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"morphic-os/backend/internal/domain"
	"strings"
	"sync/atomic"
	"time"
)

// SleepCycleConfig defines the thresholds for pruning memories.
type SleepCycleConfig struct {
	MaxLastRecallDays  int // Days since last recall to consider for pruning
	MaxAccessFrequency int // Maximum access frequency to be considered for pruning
}

// NightlySleepCycle represents the background daemon responsible for pruning cognitive memory.
type NightlySleepCycle struct {
	repo        domain.MemoryRepository
	agent       Agent
	config      SleepCycleConfig
	prunedCount int64 // Atomic counter for pruned vectors
}

// NewNightlySleepCycle creates a new NightlySleepCycle daemon.
func NewNightlySleepCycle(repo domain.MemoryRepository, agent Agent, config SleepCycleConfig) *NightlySleepCycle {
	if config.MaxLastRecallDays == 0 {
		config.MaxLastRecallDays = 30 // default 30 days
	}
	if config.MaxAccessFrequency == 0 {
		config.MaxAccessFrequency = 5 // default 5 accesses
	}
	return &NightlySleepCycle{
		repo:   repo,
		agent:  agent,
		config: config,
	}
}

// Run executes the sleep cycle, finding prunable vectors and deleting them.
func (n *NightlySleepCycle) Run(ctx context.Context) {
	log.Println("[SleepCycle] Starting Nightly Sleep Cycle...")

	maxLastRecallTime := time.Now().AddDate(0, 0, -n.config.MaxLastRecallDays)

	vectors, err := n.repo.GetPrunableVectors(ctx, maxLastRecallTime, n.config.MaxAccessFrequency)
	if err != nil {
		log.Printf("[SleepCycle] Error getting prunable vectors: %v\n", err)
		return
	}

	log.Printf("[SleepCycle] Found %d vectors to prune.\n", len(vectors))

	prunedThisRun := 0
	for _, v := range vectors {
		prompt := fmt.Sprintf(`Evaluate this memory: "%s". Is it a critical long-term fact/preference, or is it obsolete/useless data? Respond with a JSON object containing "action": "direct_response" and "response": "KEEP" or "response": "DISCARD".`, v.Content)

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

// GetPrunedCount returns the total number of pruned vectors since startup.
func (n *NightlySleepCycle) GetPrunedCount() int {
	return int(atomic.LoadInt64(&n.prunedCount))
}
