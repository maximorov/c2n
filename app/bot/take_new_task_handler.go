package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/user"
)

const DoNotGiveLocationExecutor = SymbHart + ` Ми нікому не передаємо вашу геолокацію`

type TakeNewTaskHandlerHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *TakeNewTaskHandlerHandler) UserRole() user.Role {
	return user.Executor
}

func (s *TakeNewTaskHandlerHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, `Перш ніж отримати завдання, вам треба вказати територію, на якій ви можете допомогти.`)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, `Поділіться, будь-ласка, локацією, натиснувши кнопку `+"\n\n["+CommandGetLocationAuto+"]\n\n"+`(Працює лише на телефоні)`+SymbPhone+"\n"+`Або якщо ви хочете обрати іншу локацію, оберіть її за допомогою кнопки `+SymbClip+`, як вказано на відео нижче `+SymbLoopDown)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	s.handler.sendVideoHowSendLocation(u.Message.Chat.ID, s.keyboard)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, fmt.Sprintf("%s %s\nВони матимуть вигляд: \n`%s`", SymbWarning, GoogleSuggestion, `50.44639862968634, 30.521755358513595`))
	msg.ParseMode = `markdown`
	msg.ReplyMarkup = GoogleMapsKeyboard
	s.handler.Ans(msg)

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, DoNotGiveLocationExecutor)
	msg.ReplyMarkup = s.keyboard
	s.handler.Ans(msg)

	return true
}
