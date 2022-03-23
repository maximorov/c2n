package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandGetLocationAuto = Symbanchor + " Надати геолокацію"

type TakeLocationAutoHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *TakeLocationAutoHandler) Handle(ctx context.Context, u *tgbotapi.Update) {

}
