package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandInformation = "Довідка"

type AboutHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *AboutHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	//msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = s.keyboard
	msg.Text = Information

	s.handler.Ans(msg)
}
