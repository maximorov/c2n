package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core/db"
	"helpers/app/domains/user"
)

const SetExecutorLocation = `Зазначте де ви знаходитесь`

type AfterExecutorLocationSetHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *AfterExecutorLocationSetHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	usr := ctx.Value(`user`).(*user.User)
	err := s.registerExecutor(ctx, u, usr.ID)
	if err != nil {
		zap.S().Error(err)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = s.keyboard
	msg.Text = `Вы подписаны на информирование`
	s.handler.Ans(msg)

	msg.Text = "Оберіть територію де ви зможете допомогти"
	s.handler.Ans(msg)
}

func (s *AfterExecutorLocationSetHandler) registerExecutor(ctx context.Context, u *tgbotapi.Update, userID int) error {
	su := s.handler.ExecutorService
	_, err := su.GetOneByUserID(ctx, userID)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			pos := db.CreatePoint(u.Message.Location.Latitude, u.Message.Location.Longitude)
			err = su.CreateOne(ctx, userID, 1000, "", &pos)
			if err != nil {
				zap.S().Error(err)
			}

			return nil
		}
	} else {
		// TODO: update position
	}

	return err
}
