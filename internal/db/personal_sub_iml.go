package db

import (
	"context"
	"errors"
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

// TODO: refactor
func (p *PersonalSubDBImpl) SetPersonalSub(
	ctx context.Context,
	chatID int64,
	nickname string) error {

	var existing models.PersonalSub
	result := p.db.
		WithContext(ctx).
		Where("chat_id = ?", chatID).
		First(&existing)

	if result.Error == nil {
		return p.db.
			WithContext(ctx).
			Where("chat_id = ?", chatID).
			Delete(&models.PersonalSub{}).Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newSub := models.PersonalSub{
			ChatID:    chatID,
			Nickname:  nickname,
			LastCheck: time.Now(),
		}
		return p.db.WithContext(ctx).Create(&newSub).Error
	}

	return result.Error
}
