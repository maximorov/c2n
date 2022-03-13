package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/task/activity"
	"strconv"
	"strings"
)

type Handler struct {
	BotApi              *tgbotapi.BotAPI
	TaskActivityService *activity.Service
}

func (s *Handler) Handle(u tgbotapi.Update) bool {
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
	}

	return handled
}
