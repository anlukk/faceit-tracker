package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/anlukk/faceit-tracker/internal/db/models"
	"gorm.io/gorm"
)

type SettingsDBImpl struct {
	db *gorm.DB
}

func NewSettngsDBImpl(db *gorm.DB) *SettingsDBImpl {
	return &SettingsDBImpl{db: db}
}

func (s *SettingsDBImpl) GetNotificationsEnabled(
	ctx context.Context,
	chatID int64) (bool, error) {
	var setting models.UserSettings
	err := s.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		First(&setting).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := s.db.WithContext(ctx).
			Create(&models.UserSettings{
				ChatID:               chatID,
				NotificationsEnabled: true,
			}).Error
		if err != nil {
			return false, fmt.Errorf("failed to create new settings: %w", err)
		}

		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to get notifications enabled: %w", err)
	}

	return setting.NotificationsEnabled, nil
}

func (s *SettingsDBImpl) SetNotificationsEnabled(
	ctx context.Context,
	chatID int64, enabled bool) error {
	var settings models.UserSettings
	var count int64

	err := s.db.
		WithContext(ctx).
		Model(&settings).
		Where("chat_id = ?", chatID).
		Count(&count).
		Error
	if err != nil {
		return fmt.Errorf("failed to get notifications enabled: %w", err)
	}
	if count == 0 {
		err := s.db.
			WithContext(ctx).
			Create(&models.UserSettings{
				ChatID:               chatID,
				NotificationsEnabled: enabled,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to create new settings: %w", err)
		}
	} else {
		err := s.db.
			WithContext(ctx).
			Model(&settings).
			Where("chat_id = ?", chatID).
			Update("notifications_enabled", enabled).
			Error
		if err != nil {
			return fmt.Errorf("failed to update settings: %w", err)
		}
	}

	return err
}

func (s *SettingsDBImpl) GetAllWithNotificationsEnabled(ctx context.Context) ([]int64, error) {
	var settings []models.UserSettings
	err := s.db.
		WithContext(ctx).
		Where("notifications_enabled = ?", true).
		Find(&settings).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all with notifications enabled: %w", err)
	}

	var chatIDs []int64
	for _, setting := range settings {
		chatIDs = append(chatIDs, setting.ChatID)
	}

	return chatIDs, nil
}

func (s *SettingsDBImpl) GetLanguage(
	ctx context.Context,
	chatID int64) string {
	var settings models.UserSettings
	err := s.db.
		WithContext(ctx).
		Where("chat_id = ?", chatID).
		Find(&settings.Language).
		Error
	if err != nil {
		return ""
	}

	return settings.Language
}

func (s *SettingsDBImpl) SetLanguage(
	ctx context.Context,
	chatID int64,
	language string) error {
	var settings models.UserSettings
	err := s.db.
		WithContext(ctx).
		Model(&settings).
		Where("chat_id = ?", chatID).
		Update("language", language).
		Error
	if err != nil {
		return fmt.Errorf("failed to set language: %w", err)
	}

	return nil
}
