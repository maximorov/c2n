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
	//usr := ctx.Value(`user`).(*user.User)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID

	//uId := 1
	//x, y := 1.0, 1.0
	//taskId, err := s.handler.TaskService.CreateTask(ctx, uId, x, y, u.Message.Text)
	//if err != nil {
	//	zap.S().Error(err)
	//}
	msg.ReplyMarkup = s.keyboard
	msg.Text = "`"
	//msg.Text = strconv.Itoa(taskId)
}
