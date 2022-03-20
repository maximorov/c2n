package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TakeNewTaskHandlerHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *TakeNewTaskHandlerHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)

	msg.ReplyMarkup = s.keyboard
	msg.Text = `Поділіться будь-ласка локацією кнопкою ` + "\n\n[" + CommandGetLocation + "]\n\n" + ` або якщо ви хочете обрати іншу локацію, оберить локацию за допомогою кнопки ` + SymbClip + `, як вказано на відео нижче`
	msg.Text = msg.Text + "\n...відео завантажується"

	s.handler.Ans(msg)
	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID, s.keyboard)
}
