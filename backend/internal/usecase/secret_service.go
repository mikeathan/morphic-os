package usecase

import (
	"context"
	"fmt"
	"morphic-os/backend/internal/domain"
	"morphic-os/backend/internal/infrastructure/security"
	"time"

	"github.com/google/uuid"
)

type SecretService struct {
	repo          domain.SecretRepository
	encryptionKey []byte
}

func NewSecretService(repo domain.SecretRepository, encryptionKey string) *SecretService {
	// In production, ensure key is 32 bytes
	keyBytes := []byte(encryptionKey)
	if len(keyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		keyBytes = padded
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	return &SecretService{
		repo:          repo,
		encryptionKey: keyBytes,
	}
}

func (s *SecretService) AddSecret(ctx context.Context, workspaceID, key, value string) (*domain.Secret, error) {
	encryptedValue, err := security.Encrypt([]byte(value), s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	secret := &domain.Secret{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       encryptedValue,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Save(ctx, secret); err != nil {
		return nil, fmt.Errorf("failed to save secret: %w", err)
	}

	return secret, nil
}

func (s *SecretService) GetSecretValue(ctx context.Context, workspaceID, key string) (string, error) {
	secret, err := s.repo.GetByKey(ctx, workspaceID, key)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	decryptedValue, err := security.Decrypt(secret.Value, s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return string(decryptedValue), nil
}

func (s *SecretService) ListSecrets(ctx context.Context, workspaceID string) ([]*domain.Secret, error) {
	return s.repo.List(ctx, workspaceID)
}

func (s *SecretService) DeleteSecret(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
