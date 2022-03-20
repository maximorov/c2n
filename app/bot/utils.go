package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (s *MessageHandler) sendVideoHowSendLocation(chatID int64) {
	ans := tgbotapi.NewVideo(chatID, tgbotapi.FilePath("files/video/getLocation.mp4"))
	ans.ReplyMarkup = ToMainKeyboard
	ans.Caption = "Поділіться будь-ласка локацією кнопкою \"Надіслати локацію\", або якщо ви хочете обрати іншу локацію, поділіться локацією як вказано на відео"
	_, err := s.BotApi.Send(ans)
	if err != nil {
		zap.S().Error(err)
	}
}
