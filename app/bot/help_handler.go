package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandHelp = "Допомогти"
const CommandTakeNewOrListMy = `Look my tasks, or take new one`
const CommandMyActiveTasks = "My active tasks"
const CommandTakeNewTask = "Take new task"

type HelpHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *HelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = s.keyboard
	s.handler.role = "executor"
	msg.Text = CommandTakeNewOrListMy
	s.handler.Ans(msg)
}
