package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/user"
)

const CommandNeedHelp = SymbSOS + " Попросити про допомогу"

type NeedHelpHandler struct {
	handler           *MessageHandler
	keyboard          tgbotapi.ReplyKeyboardMarkup
	keyboardHaveTasks tgbotapi.ReplyKeyboardMarkup
}

func (s *NeedHelpHandler) UserRole() user.Role {
	return user.Needy
}

func (s *NeedHelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	usr := ctx.Value(`user`).(*user.User)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, ``)
	msg.Text = `Ви можете створити нове завдання або переглянути стан своїх актуальних завдань, які ще ні ким не виконані.`
	if s.handler.TaskService.IsUserHaveUndoneTasks(ctx, usr.ID) {
		msg.ReplyMarkup = s.keyboardHaveTasks
	} else {
		msg.ReplyMarkup = s.keyboard
	}

	s.handler.Ans(msg)

	return true
}
