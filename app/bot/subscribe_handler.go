package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/user"
)

const (
	CommandUnsubscribe = SymbHide + " Відписатися від автоматичної розсилки"
	CommandSubscribe   = SymbCheckboxOn + " Підписатися на автоматичну розсилку"
)

type SubscribeHandler struct {
	handler   *MessageHandler
	keyboard  tgbotapi.ReplyKeyboardMarkup
	subscribe bool
}

func (s *SubscribeHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)

	_, err := s.handler.ExecutorService.SetSubscribeInfo(ctx, usr.ID, s.subscribe)
	if err != nil {
		zap.S().Error(err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = HeadKeyboard
	msg.ParseMode = `markdown`
	msg.Text = "Тепер ви не будете отримувати автоматичну розсилку про потребу допомоги поруч \U0001F972"

	s.handler.Ans(msg)
}
