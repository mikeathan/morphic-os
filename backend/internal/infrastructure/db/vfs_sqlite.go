package db

import (
	"context"
	"fmt"
	"morphic-os/backend/internal/domain"

	"gorm.io/gorm"
)

// SQLiteVirtualFileRepository implements domain.VirtualFileRepository using SQLite.
type SQLiteVirtualFileRepository struct {
	db *gorm.DB
}

// NewSQLiteVirtualFileRepository creates a new SQLite repository for VirtualFiles.
func NewSQLiteVirtualFileRepository(db *gorm.DB) *SQLiteVirtualFileRepository {
	return &SQLiteVirtualFileRepository{db: db}
}

// Create inserts a new VirtualFile into the database.
func (r *SQLiteVirtualFileRepository) Create(ctx context.Context, file *domain.VirtualFile) error {
	result := r.db.WithContext(ctx).Create(file)
	if result.Error != nil {
		return fmt.Errorf("failed to create virtual file: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a VirtualFile by its ID.
func (r *SQLiteVirtualFileRepository) GetByID(ctx context.Context, id string) (*domain.VirtualFile, error) {
	var file domain.VirtualFile
	result := r.db.WithContext(ctx).First(&file, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("virtual file not found")
		}
		return nil, fmt.Errorf("failed to get virtual file by id: %w", result.Error)
	}
	return &file, nil
}

// GetByPath retrieves a VirtualFile by its WorkspaceID and Path.
func (r *SQLiteVirtualFileRepository) GetByPath(ctx context.Context, workspaceID string, path string) (*domain.VirtualFile, error) {
	var file domain.VirtualFile
	result := r.db.WithContext(ctx).First(&file, "workspace_id = ? AND path = ?", workspaceID, path)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("virtual file not found")
		}
		return nil, fmt.Errorf("failed to get virtual file by path: %w", result.Error)
	}
	return &file, nil
}

// ListByWorkspace retrieves all virtual files for a given workspace.
func (r *SQLiteVirtualFileRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]*domain.VirtualFile, error) {
	var files []*domain.VirtualFile
	// Select everything except the content column to save memory
	result := r.db.WithContext(ctx).Select("id", "workspace_id", "path", "name", "is_dir", "size", "created_at", "updated_at").Where("workspace_id = ?", workspaceID).Find(&files)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list virtual files: %w", result.Error)
	}
	return files, nil
}

// Update modifies an existing VirtualFile.
func (r *SQLiteVirtualFileRepository) Update(ctx context.Context, file *domain.VirtualFile) error {
	result := r.db.WithContext(ctx).Save(file)
	if result.Error != nil {
		return fmt.Errorf("failed to update virtual file: %w", result.Error)
	}
	return nil
}

// Delete removes a VirtualFile by its ID.
func (r *SQLiteVirtualFileRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.VirtualFile{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete virtual file: %w", result.Error)
	}
	return nil
}
