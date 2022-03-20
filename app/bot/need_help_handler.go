package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/user"
)

const CommandNeedHelp = SymbLoudspeaker + " Попросити про допомогу"

type NeedHelpHandler struct {
	handler           *MessageHandler
	keyboard          tgbotapi.ReplyKeyboardMarkup
	keyboardHaveTasks tgbotapi.ReplyKeyboardMarkup
}

func (s *NeedHelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)
	if s.handler.TaskService.IsUserHaveUndoneTasks(ctx, usr.ID) {
		msg.ReplyMarkup = s.keyboardHaveTasks
	} else {
		msg.Text = `Ви можете зробити нове завдання або переглянути стан своїх актуальних завдань, які ще ні ким не виконані.`
		msg.ReplyMarkup = s.keyboard
	}
	s.handler.role = "needy"

	_, err := s.handler.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}
