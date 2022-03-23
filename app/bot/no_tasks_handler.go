package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/user"
)

const CommandNoTasks = SymbdontKnow + " Немає завдань в цьому радіусі. Ви можете змінити радіус пошуку."

type NoTasksHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *NoTasksHandler) UserRole() user.Role {
	return user.Executor
}

func (s *NoTasksHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, CommandNeedCollectLocation)
	msg.ReplyMarkup = s.keyboard

	return true
}
