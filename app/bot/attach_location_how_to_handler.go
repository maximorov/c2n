package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/core"
	"helpers/app/core/db"
	"helpers/app/domains/user"
	"helpers/app/domains/user/soc_net"
)

type AttachLocationHowToHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *AttachLocationHowToHandler) UserRole() user.Role {
	return user.Unknown
}

func (s *AttachLocationHowToHandler) Handle(ctx context.Context, u *tgbotapi.Update) error {
	usr := ctx.Value(`user`).(*user.User)

	sSoc := soc_net.NewService(db.GetPool())
	userSocNet, err := sSoc.GetOneByUserID(ctx, usr.ID)
	if err != nil {
		return err
	}

	defer func() {
		if userSocNet != nil && userSocNet.LastReceivedMessage != `` {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, userSocNet.LastReceivedMessage)
			s.handler.Ans(msg)
		}
	}()

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "\n\nВідео завантажується "+core.SymbLoading)
	s.handler.Ans(msg)

	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID, s.keyboard)

	return nil
}
