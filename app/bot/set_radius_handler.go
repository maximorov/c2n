package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const CommandRadius1 = "Радіус 1 км"
const CommandRadius3 = "Радіус 3 км"
const CommandRadius5 = "Радіус 5 км"

type SetRadiusHandler struct {
	handler   *MessageHandler
	keyboard  tgbotapi.ReplyKeyboardMarkup
	keyboard2 tgbotapi.ReplyKeyboardMarkup
	radius    int
}

func (s *SetRadiusHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID
	msg.ReplyMarkup = s.keyboard2
	msg.Text = CommandNoTasks

	_, err := s.handler.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}

	return

	//usr := ctx.Value(`user`).(*user.User)
	//e, err := s.handler.ExecutorRepo.FindOne(ctx, []string{`position`}, map[string]interface{}{
	//	`user_id`: usr.ID,
	//})
	//if err != nil {
	//	zap.S().Error(err)
	//	return
	//}
	//
	//err = s.setAreaForUser(ctx, u, usr.ID, s.radius)
	//if err != nil {
	//	zap.S().Error(err)
	//}
	//msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)
	//msg.ReplyMarkup = s.keyboard
	//
	//tasks, err := s.handler.TaskService.FindTasksInRadius(ctx, e.Position, float64(e.Area))
	//if err != nil {
	//	zap.S().Error(err)
	//	return
	//}
	//
	//if len(tasks) == 0 {
	//	msg.Text = "No tasks?"
	//	_, err := s.handler.BotApi.Send(msg)
	//	if err != nil {
	//		zap.S().Error(err)
	//	}
	//} else {
	//	for _, t := range tasks {
	//		tId := strconv.Itoa(t.ID)
	//		TasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`accept:` + tId)
	//		TasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`hide:` + tId)
	//		msg.ReplyMarkup = TasksListKeyboard
	//		msg.Text = "Task " + tId + "\n" + t.Text
	//		_, err := s.handler.BotApi.Send(msg)
	//		if err != nil {
	//			zap.S().Error(err)
	//		}
	//	}
	//}
	//return
}

func (s *SetRadiusHandler) setAreaForUser(ctx context.Context, update *tgbotapi.Update, userID, area int) error {
	ex, err := s.handler.ExecutorService.GetOneByUserID(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	_, err = s.handler.ExecutorService.UpdateOne(ctx,
		map[string]interface{}{
			`area`: area,
		}, map[string]interface{}{
			`id`: ex.ID,
		})
	if err != nil {
		zap.S().Error(err)

		return err
	}

	return nil
}
