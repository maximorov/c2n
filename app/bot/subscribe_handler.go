package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/user"
)

const (
	CommandUnsubscribe = "Відписатись від автоматичної розсилки"
	CommandSubscribe   = "Підписатись на автоматичну розсилку"
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
}
