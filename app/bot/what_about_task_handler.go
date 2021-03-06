package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/core"
	"helpers/app/domains/user"
)

const CommandFiilTaskText = "*Укажіть завдання у форматі:* \n" +
	core.SymbPerson + " ваше ім'я\n" +
	core.SymbTask + " опис завдання\n" +
	core.SymbContact + " як з вами зв'язатися."

type WhatFillTaskText struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *WhatFillTaskText) UserRole() user.Role {
	return user.Unknown
}

func (s *WhatFillTaskText) Handle(_ context.Context, u *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = s.keyboard
	msg.Text = "`"

	return nil
}
