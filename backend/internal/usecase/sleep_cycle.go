package usecase

import (
	"context"
	"log"
	"morphic-os/backend/internal/domain"
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
	config      SleepCycleConfig
	prunedCount int64 // Atomic counter for pruned vectors
}

// NewNightlySleepCycle creates a new NightlySleepCycle daemon.
func NewNightlySleepCycle(repo domain.MemoryRepository, config SleepCycleConfig) *NightlySleepCycle {
	if config.MaxLastRecallDays == 0 {
		config.MaxLastRecallDays = 30 // default 30 days
	}
	if config.MaxAccessFrequency == 0 {
		config.MaxAccessFrequency = 5 // default 5 accesses
	}
	return &NightlySleepCycle{
		repo:   repo,
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
