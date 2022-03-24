package bot

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/user"
	"helpers/app/domains/user/executor"
	"helpers/app/domains/user/soc_net"
	"helpers/app/usecase"
	"strconv"
	"unicode/utf8"
)

type MainTextHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *MainTextHandler) UserRole() user.Role {
	return user.Unknown
}

func (s *MainTextHandler) Handle(ctx context.Context, u *tgbotapi.Update) error {
	usr := ctx.Value(`user`).(*user.User)

	tsk, err := s.handler.TaskService.GetUsersLastRawTask(ctx, usr.ID)
	if err != nil {
		if core.IsRealError(err) {
			zap.S().Error(err)
		}
		return core.DefClientError
	} else {
		tuc := usecase.NewTaskUseCase(db.GetPool())
		// validation
		if utf8.RuneCountInString(u.Message.Text) > task.TaskTextLength {
			return core.NewClientError(core.SymbWarning + ` Завдання не створено: забагато текста`)
		}
		err = tuc.UpdateLastRawWithText(ctx, tsk.ID, u.Message.Text)
		if err != nil {
			zap.S().Error(err)
		}
		msg := tgbotapi.NewMessage(
			u.Message.Chat.ID,
			fmt.Sprintf("Ваше завдання #%d\nОчікуйте повідомлення протягом %d годин\nЩойно волонтер візьметься за ваше завдання, ми вам повідомимо.", tsk.ID, task.TaskDeadline),
		)
		msg.ReplyMarkup = ToMainKeyboard
		s.handler.Ans(msg)

		// Inform executors about new task in their area
		ts := task.NewService(db.GetPool())
		snr := soc_net.NewRepo(db.GetPool())

		executors, err := executor.NewRepo(db.GetPool()).FindMany( // find all with inform true
			ctx,
			[]string{`user_id`, `position`, `area`},
			sq.Eq{`inform`: true},
		)
		if core.IsRealError(err) {
			zap.S().Error(err)
		}
		if len(executors) > 0 {
			tskId := strconv.Itoa(tsk.ID)
			for _, ex := range executors {
				dist := ts.CountDistance(tsk.Position, ex.Position)
				if dist <= float64(ex.Area) {
					tsk.SetDistance(dist)
					snUser, err := snr.FindOne(ctx, []string{`soc_net_id`}, sq.Eq{`user_id`: ex.UserId})
					if err != nil {
						zap.S().Error(err)
						continue
					}
					socNetId, _ := strconv.Atoi(snUser.SocNetID)
					tsk.Text = u.Message.Text
					msg = tgbotapi.NewMessage(
						int64(socNetId),
						PrepareTaskText(tsk),
					)
					TasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`accept:` + tskId)
					TasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`hide:` + tskId)
					msg.ReplyMarkup = TasksListKeyboard
					s.handler.Ans(msg)
				}
			}
		}
	}

	return nil
}
