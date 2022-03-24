package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/core"
	"helpers/app/domains/user"
)

const DoNotGiveLocationExecutor = core.SymbHart + ` Ми нікому не передаємо вашу геолокацію` + "\n" + core.SymbLock + "Інщі користувачі також її не побачать"

type TakeNewTaskHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *TakeNewTaskHandler) UserRole() user.Role {
	return user.Executor
}

func (s *TakeNewTaskHandler) Handle(_ context.Context, u *tgbotapi.Update) error {
	text := fmt.Sprintf("*Перш ніж отримати завдання, вам треба вказати територію, на якій ви можете допомогти.*\n\n"+
		"- Ви можете поділитися локацією, натиснувши кнопку *%s* _(лише для Android, iOS)_\n"+
		"- Або якщо ви хочете обрати іншу локацію, прикріпити її, натиснувши %s _(лише для Android, iOS)_\n"+
		"- %s", CommandGetLocationAuto, core.SymbClip, core.GoogleSuggestion)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	msg.ParseMode = `markdown`
	msg.ReplyMarkup = GoogleMapsKeyboard
	s.handler.Ans(msg)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, DoNotGiveLocationExecutor)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	return nil
}
