package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/domains/user"
	"strconv"
)

const CommandMyTasks = `Мої завдання`

type ShowMyTasksHandler struct {
	handler   *MessageHandler
	keyboard  tgbotapi.ReplyKeyboardMarkup
	keyboardM tgbotapi.InlineKeyboardMarkup
}

func (s *ShowMyTasksHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)
	kb := s.keyboardM

	tasks, err := s.handler.TaskService.GetUserUndoneTasks(ctx, usr.ID)
	if err != nil && !errors.As(err, &pgx.ErrNoRows) {
		zap.S().Error(err)
		return
	}

	if len(tasks) > 0 {
		for _, tsk := range tasks {
			kb.InlineKeyboard[0][0].CallbackData = core.StrP(CancelCallback + `:` + strconv.Itoa(tsk.ID))
			msg.ReplyMarkup = kb
			msg.Text = fmt.Sprintf("Task #%d\n%s\n%s", tsk.ID, tsk.Text, tsk.Status)
			s.handler.Ans(msg)
		}
	} else {
		msg.Text = "Немає завдань про допомогу."
		msg.ReplyMarkup = s.keyboard
		s.handler.Ans(msg)
	}
}
