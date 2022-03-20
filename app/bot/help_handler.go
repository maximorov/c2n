package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandHelp = SymbStrength + " Допомогти"
const CommandTakeNewOrListMy = `Переглянути завдання в роботі, або взяти нове.`
const CommandMyActiveTasks = "Завдання в роботі"
const CommandTakeNewTask = "Узяти нове завдання"

type HelpHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *HelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	// msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = s.keyboard
	s.handler.role = "executor"
	msg.Text = CommandTakeNewOrListMy
	s.handler.Ans(msg)
}
