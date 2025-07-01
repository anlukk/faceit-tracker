package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/db/models"
	"gorm.io/gorm"
)

type SubDBImpl struct {
	db *gorm.DB
}

func NewSubDBImpl(db *gorm.DB) *SubDBImpl {
	return &SubDBImpl{
		db: db,
	}
}

func (s *SubDBImpl) Subscribe(ctx context.Context, chatID int64, playerID, nickname string) error {
	var count int64

	err := s.db.
		WithContext(ctx).
		Model(&models.Subscription{}).
		Where("chat_id = ? AND player_id = ?", chatID, playerID).
		Count(&count).
		Error
	if err != nil {
		return fmt.Errorf("SubDBImpl.Subscribe: %w", err)
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

func (s *SubDBImpl) Unsubscribe(ctx context.Context, chatID int64, playerID string) error {
	result := s.db.WithContext(ctx).
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

func (s *SubDBImpl) IsSubscribed(ctx context.Context, chatID int64, playerID string) (bool, error) {
	var count int64

	err := s.db.WithContext(ctx).
		Model(&models.Subscription{}).
		Where("chat_id = ? AND player_id = ?", chatID, playerID).
		Count(&count).
		Error
	if err != nil {
		return false, fmt.Errorf("SubDBImpl.IsSubscribed: %w", err)
	}

	return count > 0, nil
}

func (s *SubDBImpl) GetSubscribers(ctx context.Context, chatID int64) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := s.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Find(&subs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get subscribers: %w", err)
	}

	return subs, nil
}
