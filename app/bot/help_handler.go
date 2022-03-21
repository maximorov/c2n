package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/user"
)

const CommandHelp = SymbStrength + " Допомогти"
const CommandTakeNewTaskText = `Тут ви можете бачити ваші поточні завдання. а поки, тицяйте кнопку нижче ` + SymbLoopDown
const CommandListMyTasksText = `Ви можете переглянути свої поточні завдання, або знайти нове.`
const CommandMyActiveTasks = SymbClipbord + " Завдання в роботі"
const CommandTakeNewTask = SymbWork + " Узяти нове завдання"

type HelpHandler struct {
	handler         *MessageHandler
	keyboard        tgbotapi.ReplyKeyboardMarkup
	keyboardNoTasks tgbotapi.ReplyKeyboardMarkup
}

func (s *HelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)
	if s.handler.TaskService.IsExecutorHaveUndoneTasks(ctx, usr.ID) {
		msg.Text = CommandListMyTasksText
		msg.ReplyMarkup = s.keyboard
	} else {
		msg.Text = CommandTakeNewTaskText
		msg.ReplyMarkup = s.keyboardNoTasks
	}

	s.handler.role = "executor"
	s.handler.Ans(msg)
}
