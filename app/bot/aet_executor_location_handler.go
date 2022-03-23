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
	"strconv"
	"strings"
)

type SetExecutorLocationHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *SetExecutorLocationHandler) UserRole() user.Role {
	return user.Executor
}

func (s *SetExecutorLocationHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	usr := ctx.Value(`user`).(*user.User)

	var err error
	var created bool

	if u.Message.Location != nil {
		created, err = s.registerExecutor(ctx, u.Message.Location.Latitude, u.Message.Location.Longitude, usr.ID)
		if err != nil {
			zap.S().Error(err)
		}
	} else if coordsRegexp.Match([]byte(u.Message.Text)) {
		lonLat := strings.Split(u.Message.Text, `,`)
		lat, _ := strconv.ParseFloat(strings.Trim(lonLat[0], ` `), 64)
		lon, _ := strconv.ParseFloat(strings.Trim(lonLat[1], ` `), 64)
		created, err = s.registerExecutor(ctx, lat, lon, usr.ID)
		if err != nil {
			zap.S().Error(err)
		}
	} else {
		return false
	}

	if created {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, `Ви підписані на розсилку запитів у вашому районі.`)
		msg.ReplyMarkup = s.keyboard
		s.handler.Ans(msg)
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, SymbCompass+" Оберіть можливий радіус надання допомоги.")
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	return true
}

func (s *SetExecutorLocationHandler) registerExecutor(ctx context.Context, lat, lon float64, userID int) (bool, error) {
	var executorCreated bool

	su := s.handler.ExecutorService
	_, err := su.GetOneByUserID(ctx, userID)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			pos := db.CreatePoint(lat, lon)
			err = su.CreateOne(ctx, userID, 1000, "", &pos)
			if err != nil {
				zap.S().Error(err)
			} else {
				executorCreated = true
			}
		}
	} else {
		pos := db.CreatePoint(lat, lon)
		_, err = su.UpdateOne(
			ctx,
			map[string]interface{}{`position`: &pos},
			sq.Eq{`user_id`: userID},
		)
		if err != nil {
			zap.S().Error(err)
		}
	}

	return executorCreated, err
}
