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
	msg.Text = `Перш ніж отримати завдання, вам треба вказати територію, на якій ви можете допомогти. Поділіться, будь-ласка, локацією, натиснувши кнопку ` + "\n\n[" + CommandGetLocation + "]\n\n" + `Або якщо ви хочете обрати іншу локацію, оберіть її за допомогою кнопки ` + SymbClip + `, як вказано на відео нижче ` + SymbLoopDown

	s.handler.Ans(msg)
	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID, s.keyboard)
}
