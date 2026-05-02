package domain

import (
	"context"
	"time"
)

// VirtualFile represents a file in the Virtual File System.
type VirtualFile struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	WorkspaceID string    `json:"workspace_id" gorm:"index"`
	Path        string    `json:"path" gorm:"index"`
	Name        string    `json:"name"`
	IsDir       bool      `json:"is_dir"`
	Size        int64     `json:"size"`
	Content     []byte    `json:"content,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// VirtualFileRepository defines the interface for persisting and retrieving VirtualFiles.
type VirtualFileRepository interface {
	Create(ctx context.Context, file *VirtualFile) error
	GetByID(ctx context.Context, id string) (*VirtualFile, error)
	GetByPath(ctx context.Context, workspaceID string, path string) (*VirtualFile, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]*VirtualFile, error)
	Update(ctx context.Context, file *VirtualFile) error
	Delete(ctx context.Context, id string) error
}
