package usecase_test

import (
	"context"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/usecase"
	"testing"
	"time"
)

// MockMemoryRepository implements domain.MemoryRepository
type MockMemoryRepository struct {
	vectors      map[string]*domain.MemoryVector
	deleteCount  int
}

func NewMockMemoryRepository() *MockMemoryRepository {
	return &MockMemoryRepository{
		vectors: make(map[string]*domain.MemoryVector),
	}
}

func (m *MockMemoryRepository) Save(ctx context.Context, vector *domain.MemoryVector) error {
	m.vectors[vector.ID] = vector
	return nil
}

func (m *MockMemoryRepository) GetByID(ctx context.Context, id string) (*domain.MemoryVector, error) {
	if v, ok := m.vectors[id]; ok {
		return v, nil
	}
	return nil, nil // Or error, simplified for mock
}

func (m *MockMemoryRepository) ListAll(ctx context.Context, workspaceID string) ([]*domain.MemoryVector, error) {
	var results []*domain.MemoryVector
	for _, v := range m.vectors {
		if v.WorkspaceID == workspaceID {
			results = append(results, v)
		}
	}
	return results, nil
}

func (m *MockMemoryRepository) Delete(ctx context.Context, id string) error {
	delete(m.vectors, id)
	m.deleteCount++
	return nil
}

func (m *MockMemoryRepository) GetPrunableVectors(ctx context.Context, maxLastRecall time.Time, maxAccessFrequency int) ([]*domain.MemoryVector, error) {
	var prunable []*domain.MemoryVector
	for _, v := range m.vectors {
		if !v.CoreMemory && v.LastRecall.Before(maxLastRecall) && v.AccessFrequency < maxAccessFrequency {
			prunable = append(prunable, v)
		}
	}
	return prunable, nil
}

func TestNightlySleepCycle(t *testing.T) {
	repo := NewMockMemoryRepository()
	ctx := context.Background()

	now := time.Now()

	// Seed data
	repo.Save(ctx, &domain.MemoryVector{
		ID:              "v1",
		Content:         "Recent memory, high access",
		AccessFrequency: 10,
		LastRecall:      now.Add(-2 * 24 * time.Hour), // 2 days ago
		CoreMemory:      false,
	})

	repo.Save(ctx, &domain.MemoryVector{
		ID:              "v2",
		Content:         "Old memory, low access (should prune)",
		AccessFrequency: 1,
		LastRecall:      now.Add(-40 * 24 * time.Hour), // 40 days ago
		CoreMemory:      false,
	})

	repo.Save(ctx, &domain.MemoryVector{
		ID:              "v3",
		Content:         "Old core memory, low access (should NOT prune)",
		AccessFrequency: 1,
		LastRecall:      now.Add(-40 * 24 * time.Hour),
		CoreMemory:      true,
	})

	repo.Save(ctx, &domain.MemoryVector{
		ID:              "v4",
		Content:         "Old memory, high access (should NOT prune)",
		AccessFrequency: 20,
		LastRecall:      now.Add(-40 * 24 * time.Hour),
		CoreMemory:      false,
	})

	config := usecase.SleepCycleConfig{
		MaxLastRecallDays:  30,
		MaxAccessFrequency: 5,
	}

	daemon := usecase.NewNightlySleepCycle(repo, config)

	if daemon.GetPrunedCount() != 0 {
		t.Fatalf("Expected 0 pruned count, got %d", daemon.GetPrunedCount())
	}

	daemon.Run(ctx)

	if daemon.GetPrunedCount() != 1 {
		t.Fatalf("Expected 1 pruned count, got %d", daemon.GetPrunedCount())
	}

	if repo.deleteCount != 1 {
		t.Fatalf("Expected 1 deletion from repo, got %d", repo.deleteCount)
	}

	// Verify v2 is gone, others remain
	if _, ok := repo.vectors["v2"]; ok {
		t.Errorf("Expected v2 to be deleted, but it remains")
	}
	if _, ok := repo.vectors["v1"]; !ok {
		t.Errorf("Expected v1 to remain")
	}
	if _, ok := repo.vectors["v3"]; !ok {
		t.Errorf("Expected v3 to remain")
	}
	if _, ok := repo.vectors["v4"]; !ok {
		t.Errorf("Expected v4 to remain")
	}
}
