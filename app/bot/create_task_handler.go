package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/user"
)

const CommandCreateNewTask = SymbCreate + ` Створити завдання`
const CommandNeedCollectLocation = `We need to collect info about you`
const DoNotGiveLocationNeedy = SymbOk + ` Ми нікому не передаємо вашу геолокацію`

type CreateTaskHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *CreateTaskHandler) UserRole() user.Role {
	return user.Needy
}

func (s *CreateTaskHandler) Handle(_ context.Context, u *tgbotapi.Update) bool {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, `Перш ніж отримати допомогу, вам треба вказати, де ви знаходитесь. Поділіться, будь ласка, локацією, натиснувши кнопку `+"\n\n["+CommandGetLocationAuto+"]\n\n"+`Або якщо ви хочете обрати іншу локацію, оберіть її за допомогою кнопки `+SymbClip+`, як вказано на відео нижче `+SymbLoopDown)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID, s.keyboard)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, DoNotGiveLocationNeedy)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	return true
}
