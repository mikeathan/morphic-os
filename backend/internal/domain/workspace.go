package domain

import (
	"context"
	"time"
)

// Workspace represents an isolated environment for the agent.
type Workspace struct {
	ID        string            `json:"id" gorm:"primaryKey"`
	Name      string            `json:"name"`
	EnvVars   map[string]string `json:"env_vars" gorm:"serializer:json"` // Serialized to JSON in the database
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// WorkspaceRepository defines the interface for persisting and retrieving Workspaces.
type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *Workspace) error
	GetByID(ctx context.Context, id string) (*Workspace, error)
	List(ctx context.Context) ([]*Workspace, error)
	Update(ctx context.Context, workspace *Workspace) error
	Delete(ctx context.Context, id string) error
}
