package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/task/activity"
	"strconv"
	"strings"
)

type CallbackHandler struct {
	BotApi              *tgbotapi.BotAPI
	TaskActivityService *activity.Service
}

func (s *CallbackHandler) Handle(u tgbotapi.Update) bool {
	handled := false

	switch {
	case strings.Contains(u.CallbackData(), `hide:`) || strings.Contains(u.CallbackData(), `accept:`):
		handled = true
		parsed := strings.Split(u.CallbackData(), `:`)
		action := parsed[0]
		taskId, _ := strconv.Atoi(parsed[1])

		switch action {
		case `accept`:
			err := s.TaskActivityService.CreateActivity(context.TODO(), 1, taskId, `taken`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Завдання треба виконати за добу`)
			//msg.ReplyToMessageID = u.Message.MessageID
			_, err = s.BotApi.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}
		case `hide`:
			err := s.TaskActivityService.CreateActivity(context.TODO(), 1, taskId, `hidden`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Це завдання візьме хтось інший`)
			//msg.ReplyToMessageID = u.Message.MessageID
			_, err = s.BotApi.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}
		}
	case strings.Contains(u.CallbackData(), `complete:`) || strings.Contains(u.CallbackData(), `refuse:`):
		handled = true
		parsed := strings.Split(u.CallbackData(), `:`)
		action := parsed[0]
		taskId, _ := strconv.Atoi(parsed[1])

		switch action {
		case `complete`:
			err := s.TaskActivityService.UpdateActivity(context.TODO(), 1, taskId, `completed`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Дякую`)
			_, err = s.BotApi.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}
		case `refuse`:
			err := s.TaskActivityService.UpdateActivity(context.TODO(), 1, taskId, `refused`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Візьме хтось інший, чи ні`)
			_, err = s.BotApi.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}
		}
	}

	return handled
}
