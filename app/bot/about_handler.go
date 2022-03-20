package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CommandInformation = SymbInfo + " Довідка"

const T1 = `Логика работы бота`
const d1 = "1. В системе регистрируются волонтеры, оставляют свою геолокацию и указывают радиус, в котором могут работать\n" +
	"2. Кому нужна помощь - описывают свою проблему и оставляют геолокацию\n" +
	"3. Волонтеры, в радиусе покрытия которых есть задание, уведомляются через бота\n" +
	"4. Волонтер принимает задачу и связывается с нуждающимся\n" +
	"5. По выполнению задания волонтер отмечает его как __выполнено__"
const T2 = `Безопасность`
const d2 = `- ` + BeCareful + "\n" +
	"- Мы храним только геолокацию, которую вы предоставили, и идентификатор вашего чата с нашим ботом \n" +
	"- Мы не показываем и не передаем ваши данные третьим лицам"
const BeCareful = `Будьте бдительны. Мы не проверяем достоверность информации, которую предоставили пользователи при работе с нашим ботом`

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
