package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const CommandNeedHelp = "Попросити допомогу"

type NeedHelpHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *NeedHelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	//usr := ctx.Value(`user`).(*user.User)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)

	//if s.handler.TaskService.IsUserHaveUndoneTasks(ctx, usr.ID) {
	msg.Text = `What exactly you want`
	msg.ReplyMarkup = s.keyboard
	//}
	s.handler.role = "needy"

	_, err := s.handler.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}
