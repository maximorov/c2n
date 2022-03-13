package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/usecase"
	"log"
	"strconv"
)

func botHandlers(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	u tgbotapi.UpdateConfig,
	connPool db.Conn,
) {
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		// Extract the command from the Message.
		switch update.Message.Text {
		case "/start":
			msg.ReplyMarkup = HeadKeyboard
		case CommandHelp:
			s := usecase.NewTaskUseCase(connPool)
			tasks, err := s.GetTasksForUser(ctx, 1)
			if err != nil {
				zap.S().Error(err)
			}
			msg.ReplyMarkup = nil
			if len(tasks) == 0 {
				msg.Text = "No tasks?"
				_, err := bot.Send(msg)
				if err != nil {
					zap.S().Error(err)
				}
			} else {
				for _, t := range tasks {
					msg.Text = t.Text
					_, err := bot.Send(msg)
					if err != nil {
						zap.S().Error(err)
					}
				}
			}
			continue
		case CommandNeedHelp:
			msg.ReplyMarkup = NeedHelpKeyboard
			//msg.Text =
		case CommandInformation:
			msg.Text = Information
		case CommandRadius1:
			msg.ReplyMarkup = GetLocationKeyboard
			//msg.Text =
		case CommandRadius3:
			msg.ReplyMarkup = GetLocationKeyboard
			//msg.Text =
		case CommandRadius5:
			msg.ReplyMarkup = GetLocationKeyboard
			//msg.Text =
		case CommandAllCity:
			//msg.ReplyMarkup =
			//msg.Text =
		case CommandChooseCity:
		//msg.ReplyMarkup =
		//msg.Text =
		case CommandTakeLocationManual:
		//msg.ReplyMarkup =
		//msg.Text =
		case CommandContinueHelp:
			//msg.ReplyMarkup =
			//msg.Text =
		case CommandToMain:
			msg.ReplyMarkup = HeadKeyboard
		case CommandNewTask:
			msg.ReplyMarkup = GetContactsKeyboard
			msg.Text = "Поділіться будь-ласка контактами, щоб з вами могли звʼязатись"
			//TODO: нормально пофиксить эту багу..
		//case CommandGetContact:
		//	msg.ReplyMarkup = GetLocationKeyboard
		//	msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
		case CommandGetLocation:
			//msg.ReplyMarkup =
		case CommandProcessHelp:
			//msg.ReplyMarkup =
			//msg.Text =
		case ``:
		case CommandCreateTask:
			s := task.NewService(connPool)
			uId := 1
			x, y := 1.0, 1.0
			taskId, err := s.CreateTask(ctx, uId, x, y, update.Message.Text)
			if err != nil {
				zap.S().Error(err)
			}
			//msg.ReplyMarkup =
			msg.Text = strconv.Itoa(taskId)
		default:

		}

		if update.Message.Contact != nil {
			msg.ReplyMarkup = GetLocationKeyboard
			msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
		}
		if update.Message.Location != nil {
			msg.ReplyMarkup = GetLocationKeyboard
			msg.Text = "Ви поділилися локацією..." + fmt.Sprintf("\nШирота: %v\nДовгота:%v", update.Message.Location.Longitude, update.Message.Location.Latitude)
		}

		_, err := bot.Send(msg)
		if err != nil {
			zap.S().Error(err)
		}
	}
}
