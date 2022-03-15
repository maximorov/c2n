package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core/db"
)

const CommandHelp = "Допомогти"

type HelpHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *HelpHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	//usr := ctx.Value(`user`).(*user.User)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID

	msg.ReplyMarkup = s.keyboard
	//_, err := s.registerExecutor(ctx, u, usr.ID)
	//if err != nil {
	//	zap.S().Error(err)
	//}
	msg.Text = SetExecutorLocation
	_, err := s.handler.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *HelpHandler) registerExecutor(ctx context.Context, update *tgbotapi.Update, userID int) (int, error) {
	su := s.handler.ExecutorService
	ex, err := su.GetOneByUserID(ctx, userID)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			pos := db.CreatePoint(0, 0)
			userExecutorID, err := su.CreateOne(ctx, userID, 0, "", pos)
			if err != nil {
				zap.S().Error(err)
			}

			return userExecutorID, nil
		}
	}

	return ex.ID, err
}
