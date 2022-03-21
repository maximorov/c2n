package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandStart = "/start"
const CommandToMain = SymbBegining + " До початку"

const BotTitle = `Тут допомагають`
const HelloText = `Бот був створений волонтерами на допомогу волонтерам, і тим, кому потрібна допомога волонтерів.`

type StartHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *StartHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = s.keyboard
	msg.ParseMode = `markdown`
	msg.Text = fmt.Sprintf("*%s*\n\n%s\n\n%s", SymbSmile+` `+BotTitle, SymbHello+` `+HelloText, SymbWarning+` `+BeCareful)

	s.handler.Ans(msg)
}
