package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/domains/user"
	"strconv"
)

const CommandRadius1 = "Радіус " + core.Symb1 + " км"
const CommandRadius3 = "Радіус " + core.Symb3 + " км"
const CommandRadius5 = "Радіус " + core.Symb5 + " км"
const CommandRadius10 = "Радіус " + core.Symb1 + core.Symb0 + " км"

type SetRadiusHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
	radius   int
}

func (s *SetRadiusHandler) UserRole() user.Role {
	return user.Executor
}

func (s *SetRadiusHandler) Handle(ctx context.Context, u *tgbotapi.Update) error {
	usr := ctx.Value(`user`).(*user.User)
	ex, err := s.handler.ExecutorService.Repo.FindOne(ctx, []string{`user_id`, `position`, `area`, `inform`},
		map[string]interface{}{
			`user_id`: usr.ID,
		})
	if err != nil {
		return err
	}
	err = s.setAreaForUser(ctx, usr.ID, s.radius)
	if err != nil {
		zap.S().Error(err)
	}

	tasks, err := s.handler.TaskService.FindTasksInRadius(ctx, ex.Position, usr.ID, float64(s.radius))
	if len(tasks) == 0 {
		// no tasks in area
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, CommandNoTasks)
		msg.ReplyMarkup = s.keyboard
		s.handler.Ans(msg)

		if ex.Inform {
			msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Якщо завдання у вашому районі з'являться – ми вам про це повідомимо у цьому чаті")
			msg.ReplyMarkup = s.keyboard
			s.handler.Ans(msg)
		}
		return nil
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)

	for _, t := range tasks {
		tId := strconv.Itoa(t.ID)
		TasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`accept:` + tId)
		TasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`hide:` + tId)
		msg.ReplyMarkup = TasksListKeyboard

		msg.Text = PrepareTaskText(t)

		s.handler.Ans(msg)
	}

	return nil
}

func (s *SetRadiusHandler) setAreaForUser(ctx context.Context, userID, area int) error {
	_, err := s.handler.ExecutorService.UpdateOne(ctx,
		map[string]interface{}{
			`area`: area,
		}, map[string]interface{}{
			`user_id`: userID,
		})
	if err != nil {
		zap.S().Error(err)

		return err
	}

	return nil
}
