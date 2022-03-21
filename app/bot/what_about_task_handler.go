package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandFiilTaskText = "*Укажіть завдання у форматі:* \n" +
	SymbPerson + " ім'я\n" +
	SymbTask + " опис завдання\n" +
	SymbContact + " як з вами зв'язатися."

type WhatFillTaskText struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *WhatFillTaskText) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = s.keyboard
	msg.Text = "`"
}
