package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/user"
)

const CommandNeedHelp = "Попросити допомогу"

type NeedHelpHandler struct {
	handler *MessageHandler
}

func (s *NeedHelpHandler) Msg() string {
	return CommandNeedHelp
}

func (s *NeedHelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)

	if s.handler.TaskService.IsUserHaveUndoneTasks(ctx, usr.ID) {
		msg.Text = `What exactly you want`
		msg.ReplyMarkup = NeedHelpHaveTasksKeyboard
	} else {
		//msg.ReplyToMessageID = u.Message.MessageID
		msg.Text = `We need to collect info about you`
		msg.ReplyMarkup = GetLocationKeyboard
	}

	_, err := s.handler.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}
