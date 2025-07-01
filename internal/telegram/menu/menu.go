package menu

import (
	"go.uber.org/zap"
	"sync"
)

type MenuState struct {
	Type      string
	MessageID int
}

type MenuManager struct {
	states sync.Map
	logger *zap.SugaredLogger
}

func NewMenuManager(logger *zap.SugaredLogger) *MenuManager {
	return &MenuManager{
		logger: logger,
	}
}

func (m *MenuManager) SetActive(chatID int64, menuType string, messageID int) {
	m.states.Store(chatID, MenuState{
		Type:      menuType,
		MessageID: messageID,
	})
}

func (m *MenuManager) GetActive(chatID int64) (MenuState, bool) {
	val, ok := m.states.Load(chatID)
	if !ok {
		return MenuState{}, false
	}
	return val.(MenuState), true
}

func (m *MenuManager) Clear(chatID int64) {
	m.states.Delete(chatID)
}
