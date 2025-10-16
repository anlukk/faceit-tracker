package notifier

type Messenger interface {
	SendMessage(chatID int64, text string) error
}
