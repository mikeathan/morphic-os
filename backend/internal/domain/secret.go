package domain

import (
	"context"
	"time"
)

type Secret struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	WorkspaceID string    `json:"workspace_id" gorm:"index"`
	Key         string    `json:"key"`
	Value       string    `json:"-"` // encrypted, never serialize
	CreatedAt   time.Time `json:"created_at"`
}

type SecretRepository interface {
	Save(ctx context.Context, secret *Secret) error
	GetByKey(ctx context.Context, workspaceID, key string) (*Secret, error)
	List(ctx context.Context, workspaceID string) ([]*Secret, error)
	Delete(ctx context.Context, id string) error
}
