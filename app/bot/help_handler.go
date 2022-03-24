package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/core"
	"helpers/app/domains/user"
)

const CommandHelp = core.SymbStrength + " Допомогти"
const CommandTakeNewTaskText = `Тут ви можете бачити ваші поточні завдання. А поки, тицяйте кнопку нижче ` + core.SymbLoopDown
const CommandListMyTasksText = `Ви можете переглянути свої поточні завдання, або знайти нове.`
const CommandMyActiveTasks = core.SymbClipbord + " Завдання в роботі"
const CommandTakeNewTask = core.SymbWork + " Узяти нове завдання"

type HelpHandler struct {
	handler         *MessageHandler
	keyboard        tgbotapi.ReplyKeyboardMarkup
	keyboardNoTasks tgbotapi.ReplyKeyboardMarkup
}

func (s *HelpHandler) UserRole() user.Role {
	return user.Executor
}

func (s *HelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) error {
	usr := ctx.Value(`user`).(*user.User)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)
	if s.handler.TaskService.IsExecutorHaveUndoneTasks(ctx, usr.ID) {
		msg.Text = CommandListMyTasksText
		msg.ReplyMarkup = s.keyboard
	} else {
		msg.Text = CommandTakeNewTaskText
		msg.ReplyMarkup = s.keyboardNoTasks
	}

	s.handler.Ans(msg)

	return nil
}
