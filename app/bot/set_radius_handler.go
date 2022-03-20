package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/domains/user"
	"strconv"
	"time"
)

const CommandRadius1 = "Радіус 1 км"
const CommandRadius3 = "Радіус 3 км"
const CommandRadius5 = "Радіус 5 км"
const CommandRadius10 = "Радіус 10 км"

type SetRadiusHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
	radius   int
}

func (s *SetRadiusHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)
	e, err := s.handler.ExecutorRepo.FindOne(ctx, []string{`user_id`, `position`, `area`},
		map[string]interface{}{
			`user_id`: usr.ID,
		})
	if err != nil {
		zap.S().Error(err)
		return
	}
	err = s.setAreaForUser(ctx, usr.ID, s.radius)
	if err != nil {
		zap.S().Error(err)
	}

	tasks, err := s.handler.TaskService.FindTasksInRadius(ctx, e.Position, usr.ID, float64(e.Area))
	if len(tasks) == 0 {
		// no tasks in area
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
		// msg.ReplyToMessageID = u.Message.MessageID
		msg.ReplyMarkup = s.keyboard
		msg.Text = CommandNoTasks

		s.handler.Ans(msg)
		return
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)
	msg.ReplyMarkup = s.keyboard

	for _, t := range tasks {
		tId := strconv.Itoa(t.ID)
		TasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`accept:` + tId)
		TasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`hide:` + tId)
		msg.ReplyMarkup = TasksListKeyboard

		past := time.Since(t.Created)
		hoursAgo := past.Hours()
		var pastText string
		if hoursAgo < 1 {
			pastText = `менш ніж годину тому`
		} else {
			pastText = strconv.Itoa(int(hoursAgo)) + ` годин тому`
		}

		taskText := `Завдання #` + tId + "\n"
		taskText = taskText + `Створено ` + pastText + "\n"
		taskText = taskText + t.Text
		msg.Text = taskText

		s.handler.Ans(msg)
	}

	return
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
