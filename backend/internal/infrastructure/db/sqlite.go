package db

import (
	"context"
	"fmt"
	"morphic-os/backend/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SQLiteToolRepository implements domain.ToolRepository using SQLite.
type SQLiteToolRepository struct {
	db *gorm.DB
}

// GetDB returns the underlying GORM database instance.
func (r *SQLiteToolRepository) GetDB() *gorm.DB {
	return r.db
}

// NewSQLiteToolRepository creates a new SQLite repository and runs migrations.
func NewSQLiteToolRepository(dbPath string) (*SQLiteToolRepository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite database: %w", err)
	}

	// Auto-migrate the Tool schema
	if err := db.AutoMigrate(&domain.Workspace{}, &domain.Tool{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &SQLiteToolRepository{db: db}, nil
}

// Create inserts a new Tool into the database.
func (r *SQLiteToolRepository) Create(ctx context.Context, tool *domain.Tool) error {
	result := r.db.WithContext(ctx).Create(tool)
	if result.Error != nil {
		return fmt.Errorf("failed to create tool: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a Tool by its ID.
func (r *SQLiteToolRepository) GetByID(ctx context.Context, id string) (*domain.Tool, error) {
	var tool domain.Tool
	result := r.db.WithContext(ctx).First(&tool, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tool not found")
		}
		return nil, fmt.Errorf("failed to get tool by id: %w", result.Error)
	}
	return &tool, nil
}

// GetByName retrieves a Tool by its Name and WorkspaceID.
func (r *SQLiteToolRepository) GetByName(ctx context.Context, workspaceID string, name string) (*domain.Tool, error) {
	var tool domain.Tool
	result := r.db.WithContext(ctx).First(&tool, "workspace_id = ? AND name = ?", workspaceID, name)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tool not found")
		}
		return nil, fmt.Errorf("failed to get tool by name: %w", result.Error)
	}
	return &tool, nil
}

// ListActive retrieves all tools where Active is true for a given workspace.
func (r *SQLiteToolRepository) ListActive(ctx context.Context, workspaceID string) ([]*domain.Tool, error) {
	var tools []*domain.Tool
	result := r.db.WithContext(ctx).Where("workspace_id = ? AND active = ?", workspaceID, true).Find(&tools)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list active tools: %w", result.Error)
	}
	return tools, nil
}

// Update modifies an existing Tool.
func (r *SQLiteToolRepository) Update(ctx context.Context, tool *domain.Tool) error {
	result := r.db.WithContext(ctx).Save(tool)
	if result.Error != nil {
		return fmt.Errorf("failed to update tool: %w", result.Error)
	}
	return nil
}

// Delete removes a Tool by its ID.
func (r *SQLiteToolRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Tool{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete tool: %w", result.Error)
	}
	return nil
}
