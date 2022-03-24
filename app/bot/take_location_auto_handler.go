package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/core"
)

const CommandGetLocationAuto = core.Symbanchor + " Надати геолокацію"

type TakeLocationAutoHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *TakeLocationAutoHandler) Handle(_ context.Context, _ *tgbotapi.Update) error {
	return nil
}
