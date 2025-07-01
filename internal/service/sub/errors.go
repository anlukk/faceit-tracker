package sub

import "errors"

var (
	ErrInvalidChatID   = errors.New("invalid chatID")
	ErrInvalidPlayerID = errors.New("invalid playerID")
	ErrInvalidNickname = errors.New("invalid nickname")
	ErrAlreadyExists   = errors.New("already exists")
)
