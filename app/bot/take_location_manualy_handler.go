package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/user"
)

const CommandGetLocationManual = SymbHand + ` Ввести геолокацію вручну(за допомогою google maps)`

type TakeLocationManualyHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *TakeLocationManualyHandler) UserRole() user.Role {
	return user.Executor
}

func (s *TakeLocationManualyHandler) Handle(_ context.Context, u *tgbotapi.Update) bool {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	// msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = s.keyboard
	msg.Text = CommandGetLocationManual

	s.handler.Ans(msg)

	return true
}
