package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/domains/user"
	"strconv"
)

const NoUndoneTasksMessage = `No undone tasks`

type MyActiveTasksHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *MyActiveTasksHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)
	tasks, err := s.handler.TaskUseCase.GetExecutorUndoneTasks(ctx, usr.ID)
	if err != nil && !errors.As(err, &pgx.ErrNoRows) {
		zap.S().Error(err)
	}
	if len(tasks) == 0 {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
		msg.ReplyToMessageID = u.Message.MessageID
		msg.ReplyMarkup = s.keyboard
		msg.Text = NoUndoneTasksMessage
		s.handler.Ans(msg)
		return
	}

	for _, t := range tasks {
		tId := strconv.Itoa(t.ID)
		ExecutorTasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`complete:` + tId)
		ExecutorTasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`refuse:` + tId)
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
		msg.ReplyToMessageID = u.Message.MessageID
		msg.ReplyMarkup = ExecutorTasksListKeyboard
		msg.Text = "Task " + tId + "\n" + t.Text
		_, err := s.handler.BotApi.Send(msg)
		if err != nil {
			zap.S().Error(err)
		}
	}
}
