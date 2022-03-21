package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandInformation = SymbInfo + " Довідка"

const T1 = SymbEyes + ` Логіка роботи бота`
const d1 = "1. У системі реєструються волонтери, надають свою геолокацію та вказують свій можливий радіус надання допомоги.\n" +
	"2. Люди, що потребують допомоги, описують свою проблему та надають свою геолокацію.\n" +
	"3. Волонтери отримують повідомлення через бот, якщо в радіусі їх дії з'являється завдання.\n" +
	"4. Волонтер приймає завдання і виходить на зв'язок з людиною, що потребує допомоги.\n" +
	"5. Коли допомогу надано, волонтер позначає завдання як \"виконане\"."
const T2 = SymbSecurity + ` Безпека`
const d2 = `- ` + BeCareful + "\n" +
	"- Ми зберігаємо лише надану вами геолокацію та ідентифікатор вашого чату з нашим ботом.\n" +
	"- Ми не показуємо та не передаємо ваші дані третім особам."
const BeCareful = `Будьте уважними. Ми не перевіряємо правдивість інформації, наданої користувачами під час роботи з нашим ботом.`

type AboutHandler struct {
	handler  *MessageHandler
	keyboard tgbotapi.ReplyKeyboardMarkup
}

func (s *AboutHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.Text = fmt.Sprintf("*%s*\n%s\n\n*%s*\n%s", T1, d1, T2, d2)
	msg.ReplyMarkup = s.keyboard
	msg.ParseMode = `markdown`

	s.handler.Ans(msg)
}
