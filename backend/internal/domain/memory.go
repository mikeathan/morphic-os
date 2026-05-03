package domain

import (
	"context"
	"time"
)

// MemoryVector represents a consolidated fact or preference in the KnowledgeBase.
// This is the metadata structure used for pruning; the actual embedding is theoretically handled by sqlite-vss.
type MemoryVector struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	WorkspaceID     string    `json:"workspace_id" gorm:"index"`
	Content         string    `json:"content"`          // The actual fact/text
	AccessFrequency int       `json:"access_frequency"` // Number of times recalled
	LastRecall      time.Time `json:"last_recall"`      // Last time this memory was accessed
	CoreMemory      bool      `json:"core_memory"`      // If true, never delete
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MemoryRepository defines the interface for interacting with the cognitive memory storage.
type MemoryRepository interface {
	Save(ctx context.Context, vector *MemoryVector) error
	GetByID(ctx context.Context, id string) (*MemoryVector, error)
	ListAll(ctx context.Context, workspaceID string) ([]*MemoryVector, error)
	Delete(ctx context.Context, id string) error

	// GetPrunableVectors returns vectors that meet the criteria for being forgotten.
	// For instance, last recall was long ago, and access frequency is low.
	GetPrunableVectors(ctx context.Context, maxLastRecall time.Time, maxAccessFrequency int) ([]*MemoryVector, error)
}
