package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *MessageHandler) sendVideoHowSendLocation(chatID int64, kb tgbotapi.ReplyKeyboardMarkup) {
	ans := tgbotapi.NewVideo(chatID, tgbotapi.FilePath("files/video/getLocation.mp4"))
	ans.ReplyMarkup = kb
	//ans.Caption = ""

	s.Ans(ans)
}

func (s *MessageHandler) sendVideoHowToHelp(chatID int64, kb tgbotapi.ReplyKeyboardMarkup) {
	ans := tgbotapi.NewVideo(chatID, tgbotapi.FilePath("files/video/howToHelp.mp4"))
	ans.ReplyMarkup = kb
	//ans.Caption = ""

	s.Ans(ans)
}

func (s *MessageHandler) sendVideoHowToGetHelp(chatID int64, kb tgbotapi.ReplyKeyboardMarkup) {
	ans := tgbotapi.NewVideo(chatID, tgbotapi.FilePath("files/video/howToGetHelp.mp4"))
	ans.ReplyMarkup = kb
	//ans.Caption = ""

	s.Ans(ans)
}
