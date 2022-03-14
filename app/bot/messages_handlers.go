package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/task"
)

func NewMessageHandler(BotApi *tgbotapi.BotAPI,
	TaskService *task.Service) *MessageHandler {
	res := &MessageHandler{BotApi: BotApi, TaskService: TaskService}
	res.Init()

	return res
}

type Handler interface {
	Handle(context.Context, *tgbotapi.Update)
}

type MessageHandler struct {
	handlers    map[string]Handler
	BotApi      *tgbotapi.BotAPI
	TaskService *task.Service
}

func (s *MessageHandler) Init() {
	s.handlers = map[string]Handler{
		CommandNeedHelp: &NeedHelpHandler{s},
	}
}

func (s *MessageHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	if h, ok := s.handlers[u.Message.Text]; ok {
		h.Handle(ctx, u)
		return true
	}

	return false
}
