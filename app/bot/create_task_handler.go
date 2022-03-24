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
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, `Перш ніж отримати допомогу, вам треба вказати, де ви знаходитесь. Поділіться, будь ласка, локацією, натиснувши кнопку `+"\n\n["+CommandGetLocationAuto+"]\n\n"+`Або якщо ви хочете обрати іншу локацію, оберіть її за допомогою кнопки `+core.SymbClip /*+`, як вказано на відео нижче `+SymbLoopDown*/)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID, s.keyboard)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, fmt.Sprintf("%s %s\nВони матимуть вигляд: \n`%s`", core.SymbWarning, core.GoogleSuggestion, `50.44639862968634, 30.521755358513595`))
	msg.ParseMode = `markdown`
	msg.ReplyMarkup = GoogleMapsKeyboard
	s.handler.Ans(msg)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, DoNotGiveLocationNeedy)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	return nil
}
