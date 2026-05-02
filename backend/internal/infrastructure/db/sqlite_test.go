package db_test

import (
	"context"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/db"
	"testing"
)

func TestSQLiteToolRepository(t *testing.T) {
	// Use in-memory SQLite for testing
	repo, err := db.NewSQLiteToolRepository("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}

	ctx := context.Background()

	// Test Create
	tool := &domain.Tool{
		ID:          "1",
		Name:        "test_tool",
		Description: "A test tool",
		Active:      true,
	}

	err = repo.Create(ctx, tool)
	if err != nil {
		t.Fatalf("failed to create tool: %v", err)
	}

	// Test GetByID
	fetchedTool, err := repo.GetByID(ctx, "1")
	if err != nil {
		t.Fatalf("failed to get tool by id: %v", err)
	}
	if fetchedTool.Name != "test_tool" {
		t.Errorf("expected tool name 'test_tool', got '%s'", fetchedTool.Name)
	}

	// Test ListActive
	activeTools, err := repo.ListActive(ctx, "")
	if err != nil {
		t.Fatalf("failed to list active tools: %v", err)
	}
	if len(activeTools) != 1 {
		t.Errorf("expected 1 active tool, got %d", len(activeTools))
	}
}
