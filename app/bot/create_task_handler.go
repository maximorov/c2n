package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/core"
	"helpers/app/domains/user"
)

const CommandCreateNewTask = core.SymbCreate + ` Створити завдання`
const CommandNeedCollectLocation = `We need to collect info about you`
const DoNotGiveLocationNeedy = core.SymbOk + ` Ми нікому не передаємо вашу геолокацію` + "\n" + core.SymbLock + "Інщі користувачі також її не побачать"

type CreateTaskHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *CreateTaskHandler) UserRole() user.Role {
	return user.Needy
}

func (s *CreateTaskHandler) Handle(_ context.Context, u *tgbotapi.Update) error {
	text := fmt.Sprintf("*Перш ніж отримати допомогу, вам треба вказати, де ви знаходитесь.*\n\n"+
		"- Ви можете поділитися локацією, натиснувши кнопку *%s* _(лише для Android, iOS)_\n"+
		"- Або якщо ви хочете обрати іншу локацію, прикріпити її, натиснувши %s _(лише для Android, iOS)_\n"+
		"- %s", CommandGetLocationAuto, core.SymbClip, core.GoogleSuggestion)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = `markdown`
	msg.ReplyMarkup = GoogleMapsKeyboard
	s.handler.Ans(msg)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, DoNotGiveLocationNeedy)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	return nil
}
