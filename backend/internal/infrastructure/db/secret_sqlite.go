package db

import (
	"context"
	"morphic-os/backend/internal/domain"

	"gorm.io/gorm"
)

type SQLiteSecretRepository struct {
	db *gorm.DB
}

func NewSQLiteSecretRepository(db *gorm.DB) *SQLiteSecretRepository {
	return &SQLiteSecretRepository{db: db}
}

func (r *SQLiteSecretRepository) Save(ctx context.Context, secret *domain.Secret) error {
	return r.db.WithContext(ctx).Save(secret).Error
}

func (r *SQLiteSecretRepository) GetByKey(ctx context.Context, workspaceID, key string) (*domain.Secret, error) {
	var secret domain.Secret
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND key = ?", workspaceID, key).First(&secret).Error
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

func (r *SQLiteSecretRepository) List(ctx context.Context, workspaceID string) ([]*domain.Secret, error) {
	var secrets []*domain.Secret
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&secrets).Error
	return secrets, err
}

func (r *SQLiteSecretRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Secret{}, "id = ?", id).Error
}
