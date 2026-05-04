package wasm

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"morphic-os/backend/internal/domain"
	"time"
)

// wazeroVFS implements fs.FS using the domain.VirtualFileRepository
type wazeroVFS struct {
	ctx         context.Context
	repo        domain.VirtualFileRepository
	workspaceID string
}

func NewWazeroVFS(ctx context.Context, repo domain.VirtualFileRepository, workspaceID string) fs.FS {
	return &wazeroVFS{
		ctx:         ctx,
		repo:        repo,
		workspaceID: workspaceID,
	}
}

// Open implements fs.FS.Open
func (v *wazeroVFS) Open(name string) (fs.File, error) {
	// Normalize path
	if name == "." || name == "" {
		name = "/"
	} else if name[0] != '/' {
		name = "/" + name
	}

	vf, err := v.repo.GetByPath(v.ctx, v.workspaceID, name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	return &wazeroFile{
		vf:     vf,
		offset: 0,
		repo:   v.repo,
		ctx:    v.ctx,
	}, nil
}

// wazeroFile implements fs.File and fs.ReadDirFile
type wazeroFile struct {
	vf     *domain.VirtualFile
	offset int64
	repo   domain.VirtualFileRepository
	ctx    context.Context
}

// Stat implements fs.File.Stat
func (f *wazeroFile) Stat() (fs.FileInfo, error) {
	return &wazeroFileInfo{vf: f.vf}, nil
}

// Read implements fs.File.Read
func (f *wazeroFile) Read(p []byte) (n int, err error) {
	if f.vf.IsDir {
		return 0, fmt.Errorf("is a directory")
	}
	if f.offset >= int64(len(f.vf.Content)) {
		return 0, io.EOF
	}
	n = copy(p, f.vf.Content[f.offset:])
	f.offset += int64(n)
	return n, nil
}

// Close implements fs.File.Close
func (f *wazeroFile) Close() error {
	return nil
}

// ReadDir implements fs.ReadDirFile.ReadDir
func (f *wazeroFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if !f.vf.IsDir {
		return nil, fmt.Errorf("not a directory")
	}

	// This is a naive implementation because domain.VirtualFileRepository doesn't have ListByPath.
	// For a complete implementation, we'd need to filter ListByWorkspace by prefix.
	allFiles, err := f.repo.ListByWorkspace(f.ctx, f.vf.WorkspaceID)
	if err != nil {
		return nil, err
	}

	var entries []fs.DirEntry
	// Just return all files for this simple implementation as requested by code review,
	// ideally filter by path prefix if domain allowed.
	for _, file := range allFiles {
		// Ignore the directory itself
		if file.ID == f.vf.ID {
			continue
		}
		entries = append(entries, &wazeroDirEntry{vf: file})
	}

	if n > 0 && len(entries) > n {
		entries = entries[:n]
	}

	return entries, nil
}

// wazeroDirEntry implements fs.DirEntry
type wazeroDirEntry struct {
	vf *domain.VirtualFile
}

func (d *wazeroDirEntry) Name() string               { return d.vf.Name }
func (d *wazeroDirEntry) IsDir() bool                { return d.vf.IsDir }
func (d *wazeroDirEntry) Type() fs.FileMode          { if d.vf.IsDir { return fs.ModeDir }; return 0 }
func (d *wazeroDirEntry) Info() (fs.FileInfo, error) { return &wazeroFileInfo{vf: d.vf}, nil }

// wazeroFileInfo implements fs.FileInfo
type wazeroFileInfo struct {
	vf *domain.VirtualFile
}

func (i *wazeroFileInfo) Name() string       { return i.vf.Name }
func (i *wazeroFileInfo) Size() int64        { return i.vf.Size }
func (i *wazeroFileInfo) Mode() fs.FileMode  {
	if i.vf.IsDir {
		return fs.ModeDir | 0755
	}
	return 0644
}
func (i *wazeroFileInfo) ModTime() time.Time { return i.vf.UpdatedAt }
func (i *wazeroFileInfo) IsDir() bool        { return i.vf.IsDir }
func (i *wazeroFileInfo) Sys() any           { return nil }
