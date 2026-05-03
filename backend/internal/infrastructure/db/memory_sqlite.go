package db

import (
	"context"
	"fmt"
	"morphic-os/backend/internal/domain"
	"time"

	"gorm.io/gorm"
)

// SQLiteMemoryRepository implements domain.MemoryRepository using SQLite.
type SQLiteMemoryRepository struct {
	db *gorm.DB
}

// NewSQLiteMemoryRepository creates a new SQLite repository for MemoryVectors.
func NewSQLiteMemoryRepository(db *gorm.DB) *SQLiteMemoryRepository {
	return &SQLiteMemoryRepository{db: db}
}

// Save inserts or updates a MemoryVector in the database.
func (r *SQLiteMemoryRepository) Save(ctx context.Context, vector *domain.MemoryVector) error {
	result := r.db.WithContext(ctx).Save(vector)
	if result.Error != nil {
		return fmt.Errorf("failed to save memory vector: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a MemoryVector by its ID.
func (r *SQLiteMemoryRepository) GetByID(ctx context.Context, id string) (*domain.MemoryVector, error) {
	var vector domain.MemoryVector
	result := r.db.WithContext(ctx).First(&vector, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("memory vector not found")
		}
		return nil, fmt.Errorf("failed to get memory vector by id: %w", result.Error)
	}
	return &vector, nil
}

// ListAll retrieves all MemoryVectors for a given workspace.
func (r *SQLiteMemoryRepository) ListAll(ctx context.Context, workspaceID string) ([]*domain.MemoryVector, error) {
	var vectors []*domain.MemoryVector
	result := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&vectors)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list memory vectors: %w", result.Error)
	}
	return vectors, nil
}

// Delete removes a MemoryVector by its ID.
func (r *SQLiteMemoryRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.MemoryVector{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete memory vector: %w", result.Error)
	}
	return nil
}

// GetPrunableVectors returns MemoryVectors that meet the criteria for being forgotten.
func (r *SQLiteMemoryRepository) GetPrunableVectors(ctx context.Context, maxLastRecall time.Time, maxAccessFrequency int) ([]*domain.MemoryVector, error) {
	var vectors []*domain.MemoryVector
	// Criteria: Not CoreMemory AND LastRecall is before maxLastRecall AND AccessFrequency is less than maxAccessFrequency
	result := r.db.WithContext(ctx).
		Where("core_memory = ? AND last_recall < ? AND access_frequency < ?", false, maxLastRecall, maxAccessFrequency).
		Find(&vectors)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get prunable vectors: %w", result.Error)
	}
	return vectors, nil
}
