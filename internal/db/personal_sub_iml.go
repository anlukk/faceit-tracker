package db

import (
	"context"
	"fmt"
	"time"

	"github.com/anlukk/faceit-tracker/internal/db/models"
	"gorm.io/gorm"
)

type PersonalSubDBImpl struct {
	db *gorm.DB
}

func NewPersonalSubDBImpl(db *gorm.DB) *PersonalSubDBImpl {
	return &PersonalSubDBImpl{
		db: db,
	}
}

func (p *PersonalSubDBImpl) GetPersonalSub(
	ctx context.Context,
	chatID int64) (*models.PersonalSub, error) {
	var personalSub models.PersonalSub
	result := p.db.
		WithContext(ctx).
		Where("chat_id = ?", chatID).
		First(&personalSub)
	if result.Error != nil {
		return nil, result.Error
	}
	return &personalSub, nil
}

func (p *PersonalSubDBImpl) SetPersonalSub(
	ctx context.Context,
	chatID int64) error {

	err := p.db.
		WithContext(ctx).
		Where("chat_id = ?", chatID).
		Delete(&models.PersonalSub{}).
		Error
	if err != nil {
		return fmt.Errorf("failed to delete personal sub: %w", err)
	}

	c := p.db.WithContext(ctx).Create(&models.PersonalSub{
		ChatID: chatID,
		//PlayerID:  playerID,
		LastCheck: time.Now(),
	}).Error

	return c
}
