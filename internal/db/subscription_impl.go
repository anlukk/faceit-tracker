package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/db/models"
	"gorm.io/gorm"
)

type SubscriptionDBImpl struct {
	db *gorm.DB
}

func NewSubscriptionDBImpl(db *gorm.DB) *SubscriptionDBImpl {
	return &SubscriptionDBImpl{db: db}
}

func (s *SubscriptionDBImpl) Subscribe(ctx context.Context, chatID int64, playerID, nickname string) error {
	var count int64

	err := s.db.
		WithContext(ctx).
		Model(&models.Subscription{}).
		Where("chat_id = ? AND player_id = ?", chatID, playerID).
		Count(&count).
		Error
	if err != nil {
		return fmt.Errorf("SubscriptionDBImpl.Subscribe: %w", err)
	}

	if count > 0 {
		return errors.New("already subscribed")
	}

	sub := models.Subscription{
		ChatID:   chatID,
		PlayerID: playerID,
		Nickname: nickname,
	}

	return s.db.WithContext(ctx).Create(&sub).Error
}

func (s *SubscriptionDBImpl) Unsubscribe(ctx context.Context, chatID int64, playerID string) error {
	result := s.db.
		WithContext(ctx).
		Where("chat_id = ? AND player_id = ?", chatID, playerID).
		Delete(&models.Subscription{})

	if result.Error != nil {
		return fmt.Errorf("unsubscribe failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (s *SubscriptionDBImpl) IsSubscribed(ctx context.Context, chatID int64, playerID string) (bool, error) {
	var count int64

	err := s.db.
		WithContext(ctx).
		Model(&models.Subscription{}).
		Where("chat_id = ? AND player_id = ?", chatID, playerID).
		Count(&count).
		Error
	if err != nil {
		return false, fmt.Errorf("SubscriptionDBImpl.IsSubscribed: %w", err)
	}

	return count > 0, nil
}

func (s *SubscriptionDBImpl) GetSubscribers(ctx context.Context, chatID int64) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := s.db.
		WithContext(ctx).
		Where("chat_id = ?", chatID).
		Find(&subs).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get subscribers: %w", err)
	}

	return subs, nil
}
