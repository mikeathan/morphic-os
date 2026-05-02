package db

import (
	"context"
	"fmt"
	"morphic-os/backend/internal/domain"

	"gorm.io/gorm"
)

// SQLiteWorkspaceRepository implements domain.WorkspaceRepository using SQLite.
type SQLiteWorkspaceRepository struct {
	db *gorm.DB
}

// NewSQLiteWorkspaceRepository creates a new SQLite repository for Workspaces.
func NewSQLiteWorkspaceRepository(db *gorm.DB) *SQLiteWorkspaceRepository {
	return &SQLiteWorkspaceRepository{db: db}
}

// Create inserts a new Workspace into the database.
func (r *SQLiteWorkspaceRepository) Create(ctx context.Context, workspace *domain.Workspace) error {
	result := r.db.WithContext(ctx).Create(workspace)
	if result.Error != nil {
		return fmt.Errorf("failed to create workspace: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a Workspace by its ID.
func (r *SQLiteWorkspaceRepository) GetByID(ctx context.Context, id string) (*domain.Workspace, error) {
	var workspace domain.Workspace
	result := r.db.WithContext(ctx).First(&workspace, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("workspace not found")
		}
		return nil, fmt.Errorf("failed to get workspace by id: %w", result.Error)
	}
	return &workspace, nil
}

// List retrieves all Workspaces.
func (r *SQLiteWorkspaceRepository) List(ctx context.Context) ([]*domain.Workspace, error) {
	var workspaces []*domain.Workspace
	result := r.db.WithContext(ctx).Find(&workspaces)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", result.Error)
	}
	return workspaces, nil
}

// Update modifies an existing Workspace.
func (r *SQLiteWorkspaceRepository) Update(ctx context.Context, workspace *domain.Workspace) error {
	result := r.db.WithContext(ctx).Save(workspace)
	if result.Error != nil {
		return fmt.Errorf("failed to update workspace: %w", result.Error)
	}
	return nil
}

// Delete removes a Workspace by its ID.
func (r *SQLiteWorkspaceRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Workspace{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete workspace: %w", result.Error)
	}
	return nil
}
