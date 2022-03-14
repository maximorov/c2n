package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/task"
	"helpers/app/domains/user/executor"
)

func NewMessageHandler(BotApi *tgbotapi.BotAPI,
	ts *task.Service,
	es *executor.Repository,
) *MessageHandler {
	res := &MessageHandler{BotApi: BotApi, TaskService: ts, ExecutorRepo: es}
	res.Init()

	return res
}

type Handler interface {
	Handle(context.Context, *tgbotapi.Update)
}

type MessageHandler struct {
	handlers     map[string]Handler
	BotApi       *tgbotapi.BotAPI
	TaskService  *task.Service
	ExecutorRepo *executor.Repository
}

func (s *MessageHandler) Init() {
	setRadiusH := &SetRadiusHandler{s}
	s.handlers = map[string]Handler{
		CommandNeedHelp: &NeedHelpHandler{s},
		CommandRadius1:  setRadiusH,
		CommandRadius3:  setRadiusH,
		CommandRadius5:  setRadiusH,
	}
}

func (s *MessageHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	if h, ok := s.handlers[u.Message.Text]; ok {
		h.Handle(ctx, u)
		return true
	}

	return false
}
