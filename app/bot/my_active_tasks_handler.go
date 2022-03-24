package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/domains/user"
	"strconv"
	"time"
)

const NoUndoneTasksMessage = `Немає завдань у роботі.`

type MyActiveTasksHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *MyActiveTasksHandler) UserRole() user.Role {
	return user.Executor
}

func (s *MyActiveTasksHandler) Handle(ctx context.Context, u *tgbotapi.Update) error {
	usr := ctx.Value(`user`).(*user.User)
	tasks, err := s.handler.TaskUseCase.GetExecutorUndoneTasks(ctx, usr.ID)
	if core.IsRealError(err) {
		zap.S().Error(err)
	}
	if len(tasks) == 0 {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, NoUndoneTasksMessage)
		msg.ReplyMarkup = s.keyboard
		s.handler.Ans(msg)
		return nil
	}

	for _, t := range tasks {
		tId := strconv.Itoa(t.ID)
		ExecutorTasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`complete:` + tId)
		ExecutorTasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`refuse:` + tId)

		past := time.Since(t.Created)
		hoursAgo := past.Hours()
		var pastText string
		if hoursAgo < 1 {
			pastText = `менш ніж годину тому`
		} else {
			pastText = strconv.Itoa(int(hoursAgo)) + ` годин тому`
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
		msg.ReplyMarkup = ExecutorTasksListKeyboard
		msg.Text = fmt.Sprintf("%s Завдання #%s\n\n%s\n\n%s", core.SymbTask, tId, t.Text, pastText)
		s.handler.Ans(msg)
	}

	return nil
}
