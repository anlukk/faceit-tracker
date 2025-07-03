package settings

import (
	"context"
	"errors"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/db"
)

var (
	ErrInvalidChatID = errors.New("invalid chatID")
)

type Settings interface {
	GetNotificationsEnabled(ctx context.Context, chatID int64) (bool, error)
	SetNotificationsEnabled(ctx context.Context, chatID int64, enabled bool) error
	GetAllWithNotificationsEnabled(ctx context.Context) ([]int64, error)
	GetLanguage(ctx context.Context, chatID int64) string
	SetLanguage(ctx context.Context, chatID int64, language string) error
}

type Service struct {
	repo db.SettingsDB
}

func NewService(repo db.SettingsDB) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetNotificationsEnabled(ctx context.Context, chatID int64) (bool, error) {
	if chatID < 0 {
		return false, ErrInvalidChatID
	}

	return s.repo.GetNotificationsEnabled(ctx, chatID)
}

func (s *Service) SetNotificationsEnabled(ctx context.Context, chatID int64, enabled bool) error {
	if chatID < 0 {
		return ErrInvalidChatID
	}

	err := s.repo.SetNotificationsEnabled(ctx, chatID, enabled)
	if err != nil {
		return fmt.Errorf("failed to set notifications enabled: %w", err)
	}

	return nil
}

func (s *Service) GetAllWithNotificationsEnabled(ctx context.Context) ([]int64, error) {
	return s.repo.GetAllWithNotificationsEnabled(ctx)
}

func (s *Service) GetLanguage(ctx context.Context, chatID int64) string {
	if chatID < 0 {
		return fmt.Sprintf("invalid chatID: %d", chatID)
	}

	language := s.repo.GetLanguage(ctx, chatID)
	if language == "" {
		return "en"
	}

	return language
}

func (s *Service) SetLanguage(ctx context.Context, chatID int64, language string) error {
	if chatID < 0 {
		return ErrInvalidChatID
	}

	if language == "" {
		return errors.New("language can't be empty")
	}

	err := s.repo.SetLanguage(ctx, chatID, language)
	if err != nil {
		return fmt.Errorf("failed to set language: %w", err)
	}

	return nil
}
