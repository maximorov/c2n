package bot

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core/db"
	"helpers/app/domains/user"
)

const SetExecutorLocation = `Укажіть ваше місце перебування`

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
	msg.Text = `Ви підписані на розсилку запитів у вашому районі.`
	s.handler.Ans(msg)

	msg.Text = SymbCompass + " Оберіть можливий радіус надання допомоги."
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
		pos := db.CreatePoint(u.Message.Location.Latitude, u.Message.Location.Longitude)
		_, err = su.UpdateOne(
			ctx,
			map[string]interface{}{`position`: &pos},
			sq.Eq{`user_id`: userID},
		)
		if err != nil {
			zap.S().Error(err)
		}
	}

	return err
}
