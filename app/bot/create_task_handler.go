package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandCreateNewTask = SymbCreate + ` Створити завдання`
const CommandNeedCollectLocation = `We need to collect info about you`

type CreateTaskHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *CreateTaskHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = s.keyboard
	msg.Text = `Перш ніж отримати допомогу, вам треба вказати, де ви знакодитесь. Поділіться, будь-ласка, локацією, натиснувши кнопку ` + "\n\n[" + CommandGetLocation + "]\n\n" + `Або якщо ви хочете обрати іншу локацію, оберіть її за допомогою кнопки ` + SymbClip + `, як вказано на відео нижче ` + SymbLoopDown

	s.handler.Ans(msg)
	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID, s.keyboard)
}
