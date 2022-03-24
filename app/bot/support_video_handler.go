package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/domains/user"
)

const (
	CommandSendVideoHowHelp    = "Відеоінструкція як попросити про допомогу"
	CommandSendVideoHowGetHelp = "Відеоінструкція як знайти людей, які потребують допомоги"
)

type SupportInformationHendler struct {
	handler   *MessageHandler
	keyboard  tgbotapi.ReplyKeyboardMarkup
	typeVideo string
}

func (s *SupportInformationHendler) UserRole() user.Role {
	return user.Executor
}

func (s *SupportInformationHendler) Handle(ctx context.Context, u *tgbotapi.Update) bool {

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyMarkup = s.keyboard
	msg.ParseMode = `markdown`

	if s.typeVideo == CommandSendVideoHowHelp {
		s.handler.sendVideoHowToHelp(u.Message.Chat.ID, s.keyboard)
		msg.Text = "Щоб знайти людей, які потребують допомоги зробіть так як показано на відео👆"
	}

	if s.typeVideo == CommandSendVideoHowGetHelp {
		s.handler.sendVideoHowToGetHelp(u.Message.Chat.ID, s.keyboard)
		msg.Text = "Щоб знайти людей, які потребують допомоги зробіть так як показано на відео👆"
	}

	s.handler.Ans(msg)

	return true
}
