package sub

import (
	"context"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/db"
	"github.com/anlukk/faceit-tracker/internal/db/models"
)

type Service struct {
	repo db.SubDB
}

func NewService(repo db.SubDB) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Subscribe(ctx context.Context, chatID int64, playerID, nickname string) error {
	if chatID == 0 {
		return ErrInvalidChatID
	}

	if playerID == "" {
		return ErrInvalidPlayerID
	}

	if nickname == "" {
		return ErrInvalidNickname
	}

	exists, err := s.repo.IsSubscribed(ctx, chatID, playerID)
	if err != nil {
		return fmt.Errorf("check subscription failed: %w", err)
	}

	if exists {
		return ErrAlreadyExists
	}

	if err := s.repo.Subscribe(ctx, chatID, playerID, nickname); err != nil {
		return fmt.Errorf("create subscription failed: %w", err)
	}

	return nil
}

func (s *Service) Unsubscribe(ctx context.Context, chatID int64, playerID string) error {
	if chatID < 0 {
		return ErrInvalidChatID
	}

	if playerID == "" {
		return ErrInvalidPlayerID
	}

	err := s.repo.Unsubscribe(ctx, chatID, playerID)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	return nil
}

func (s *Service) IsSubscribed(ctx context.Context, chatID int64, playerID string) (bool, error) {
	if chatID < 0 {
		return false, ErrInvalidChatID
	}

	if playerID == "" {
		return false, ErrInvalidPlayerID
	}

	isSubscribed, err := s.repo.IsSubscribed(ctx, chatID, playerID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is subscribed: %w", err)
	}

	return isSubscribed, nil
}

func (s *Service) GetSubscribers(ctx context.Context, chatID int64) ([]models.Subscription, error) {
	subs, err := s.repo.GetSubscribers(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribers: %w", err)
	}

	return subs, nil
}
