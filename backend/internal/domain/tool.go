package domain

import (
	"context"
	"time"
)

// Tool represents a synthesized microservice tool that the agent can use.
type Tool struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	WorkspaceID string    `json:"workspace_id" gorm:"uniqueIndex:idx_workspace_name"` // Link to the Workspace
	Name        string    `json:"name" gorm:"uniqueIndex:idx_workspace_name"` // Allow same tool name in different workspaces but unique per workspace ideally
	Description string    `json:"description"`
	JSONSchema  string    `json:"json_schema"` // Stored as a JSON string
	Language    string    `json:"language"`
	SourceCode  string    `json:"source_code"`
	WasmBinary  []byte    `json:"wasm_binary"` // Compiled WebAssembly binary
	ExecStats   int       `json:"exec_stats"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToolRepository defines the interface for persisting and retrieving Tools.
type ToolRepository interface {
	Create(ctx context.Context, tool *Tool) error
	GetByID(ctx context.Context, id string) (*Tool, error)
	GetByName(ctx context.Context, workspaceID string, name string) (*Tool, error)
	ListActive(ctx context.Context, workspaceID string) ([]*Tool, error)
	Update(ctx context.Context, tool *Tool) error
	Delete(ctx context.Context, id string) error
}
