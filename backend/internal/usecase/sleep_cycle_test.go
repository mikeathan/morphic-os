package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/llm"
	"morphic-os/backend/internal/usecase"
	"strings"
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

func (m *MockMemoryRepository) SearchSimilar(ctx context.Context, queryEmbedding []float32, limit int) ([]*domain.MemoryVector, error) {
	// Simple mock implementation returning an empty list or limited vectors
	var results []*domain.MemoryVector
	count := 0
	for _, v := range m.vectors {
		if count >= limit {
			break
		}
		results = append(results, v)
		count++
	}
	return results, nil
}

// MockVirtualFileRepository implements domain.VirtualFileRepository
type MockVirtualFileRepository struct {
	files map[string]*domain.VirtualFile
	deleted map[string]bool
}

func (m *MockVirtualFileRepository) Create(ctx context.Context, file *domain.VirtualFile) error { return nil }
func (m *MockVirtualFileRepository) GetByID(ctx context.Context, id string) (*domain.VirtualFile, error) { return nil, nil }
func (m *MockVirtualFileRepository) GetByPath(ctx context.Context, workspaceID string, path string) (*domain.VirtualFile, error) {
	for _, f := range m.files {
		if f.WorkspaceID == workspaceID && f.Path == path {
			return f, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockVirtualFileRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*domain.VirtualFile, error) { return nil, nil }
func (m *MockVirtualFileRepository) Update(ctx context.Context, file *domain.VirtualFile) error { return nil }
func (m *MockVirtualFileRepository) Delete(ctx context.Context, id string) error {
	m.deleted[id] = true
	return nil
}

// MockWorkspaceRepository implements domain.WorkspaceRepository
type MockWorkspaceRepository struct {
	workspaces []*domain.Workspace
}

func (m *MockWorkspaceRepository) Create(ctx context.Context, workspace *domain.Workspace) error { return nil }
func (m *MockWorkspaceRepository) GetByID(ctx context.Context, id string) (*domain.Workspace, error) { return nil, nil }
func (m *MockWorkspaceRepository) List(ctx context.Context) ([]*domain.Workspace, error) {
	return m.workspaces, nil
}
func (m *MockWorkspaceRepository) Update(ctx context.Context, workspace *domain.Workspace) error { return nil }
func (m *MockWorkspaceRepository) Delete(ctx context.Context, id string) error { return nil }

func TestNightlySleepCycle(t *testing.T) {
	repo := NewMockMemoryRepository()
	vfsRepo := &MockVirtualFileRepository{
		files: make(map[string]*domain.VirtualFile),
		deleted: make(map[string]bool),
	}
	workspaceRepo := &MockWorkspaceRepository{
		workspaces: []*domain.Workspace{{ID: "ws1"}},
	}
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

	repo.Save(ctx, &domain.MemoryVector{
		ID:              "v5",
		Content:         "Old memory, low access, but extremely useful (should KEEP via LLM)",
		AccessFrequency: 1,
		LastRecall:      now.Add(-40 * 24 * time.Hour),
		CoreMemory:      false,
	})

	repo.Save(ctx, &domain.MemoryVector{
		ID:              "v6",
		Content:         "Old memory, low access, causes agent error (should DISCARD as fallback)",
		AccessFrequency: 1,
		LastRecall:      now.Add(-40 * 24 * time.Hour),
		CoreMemory:      false,
	})

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	vfsRepo.files["f1"] = &domain.VirtualFile{
		ID: "f1",
		WorkspaceID: "ws1",
		Path: fmt.Sprintf("/var/logs/chat/%s.jsonl", yesterday),
		Content: []byte(`{"timestamp":"2023-01-01T00:00:00Z","role":"user","content":"I like pizza"}
{"timestamp":"2023-01-01T00:01:00Z","role":"assistant","content":"Noted."}`),
	}

	agent := llm.NewMockAgent()
	agent.EvaluateTaskFunc = func(ctx context.Context, task string, tools []*domain.Tool) (string, error) {
		if strings.Contains(task, "causes agent error") {
			return "", fmt.Errorf("simulated agent error")
		}

		if strings.Contains(task, "Extract all new factual information") {
			return `["User likes pizza", "System responded with Noted."]`, nil
		}

		resp := usecase.AgentResponse{
			Action: "direct_response",
		}
		if strings.Contains(task, "extremely useful") {
			resp.Response = "KEEP"
		} else {
			resp.Response = "DISCARD"
		}
		b, _ := json.Marshal(resp)
		return string(b), nil
	}

	config := usecase.SleepCycleConfig{
		MaxLastRecallDays:  30,
		MaxAccessFrequency: 5,
	}

	daemon := usecase.NewNightlySleepCycle(repo, vfsRepo, workspaceRepo, agent, config)

	if daemon.GetPrunedCount() != 0 {
		t.Fatalf("Expected 0 pruned count, got %d", daemon.GetPrunedCount())
	}

	daemon.Run(ctx)

	// Pruned count should be 2 now: v2 (discard via LLM) and v6 (discard via error fallback)
	if daemon.GetPrunedCount() != 2 {
		t.Fatalf("Expected 2 pruned count, got %d", daemon.GetPrunedCount())
	}

	if repo.deleteCount != 2 {
		t.Fatalf("Expected 2 deletion from repo, got %d", repo.deleteCount)
	}

	// Verify v2 and v6 are gone, others remain
	if _, ok := repo.vectors["v2"]; ok {
		t.Errorf("Expected v2 to be deleted, but it remains")
	}
	if _, ok := repo.vectors["v6"]; ok {
		t.Errorf("Expected v6 to be deleted, but it remains")
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

	// Verify log consolidation
	if !vfsRepo.deleted["f1"] {
		t.Errorf("Expected VFS file f1 to be deleted")
	}

	foundPizza := false
	for _, v := range repo.vectors {
		if strings.Contains(v.Content, "User likes pizza") {
			foundPizza = true
		}
	}
	if !foundPizza {
		t.Errorf("Expected consolidated fact 'User likes pizza' to be in memory repo")
	}

	// Verify v5 was kept and its LastRecall updated
	if v5, ok := repo.vectors["v5"]; !ok {
		t.Errorf("Expected v5 to remain")
	} else {
		// LastRecall should be close to now
		if time.Since(v5.LastRecall) > time.Minute {
			t.Errorf("Expected v5 LastRecall to be updated, but it is %v", v5.LastRecall)
		}
	}
}
