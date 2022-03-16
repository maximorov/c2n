package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandFiilTaskText = "Write your problem in format: contact name, issue text"

type WhatFillTaskText struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *WhatFillTaskText) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = s.keyboard
	msg.Text = "`"
}
