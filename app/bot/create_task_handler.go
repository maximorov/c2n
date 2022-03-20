package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const CommandCreateNewTask = `Створити завдання`
const CommandNeedCollectLocation = `We need to collect info about you`

type CreateTaskHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *CreateTaskHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	// msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = s.keyboard

	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID)

	_, err := s.handler.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}
