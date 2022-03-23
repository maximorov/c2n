package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/domains/user"
	"strconv"
	"strings"
)

type SetTaskLocationHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *SetTaskLocationHandler) UserRole() user.Role {
	return user.Needy
}

func (s *SetTaskLocationHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	usr := ctx.Value(`user`).(*user.User)

	if u.Message.Location != nil {
		err := s.handler.TaskUseCase.CreateRawTask(ctx, usr.ID, u.Message.Location.Latitude, u.Message.Location.Longitude)
		if err != nil {
			zap.S().Error(err)
		}
	} else if coordsRegexp.Match([]byte(u.Message.Text)) {
		lonLat := strings.Split(u.Message.Text, `,`)
		lat, _ := strconv.ParseFloat(lonLat[0], 64)
		lon, _ := strconv.ParseFloat(lonLat[1], 64)
		err := s.handler.TaskUseCase.CreateRawTask(ctx, usr.ID, lat, lon)
		if err != nil {
			zap.S().Error(err)
		}
	} else {
		return false
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, CommandFiilTaskText)
	msg.ParseMode = `markdown`
	msg.ReplyMarkup = ToMainKeyboard
	s.handler.Ans(msg)

	return true
}
