package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/task"
)

func NewMessageHandler(BotApi *tgbotapi.BotAPI,
	TaskService *task.Service) *MessageHandler {
	res := &MessageHandler{BotApi: BotApi, TaskService: TaskService}
	res.Init()

	return res
}

type Handler interface {
	Handle(*tgbotapi.Update)
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

func (s *MessageHandler) Handle(u *tgbotapi.Update) bool {
	if h, ok := s.handlers[u.Message.Text]; ok {
		h.Handle(u)
		return true
	}

	return false
}

const CommandNeedHelp = "Попросити допомогу"

type NeedHelpHandler struct {
	handler *MessageHandler
}

func (s *NeedHelpHandler) Msg() string {
	return CommandNeedHelp
}

func (s *NeedHelpHandler) Handle(u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID
	//msg.ReplyMarkup = NeedHelpKeyboard
	//msg.Text =

	_, err := s.handler.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}
