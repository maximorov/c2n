package main

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"helpers/app/bot"
	"helpers/app/core"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/task/activity"
	"helpers/app/domains/user"
	"helpers/app/domains/user/soc_net"
	"helpers/app/usecase"
	"log"
	"strconv"
	"time"
)

func botHandlers(
	ctx context.Context,
	botApi *tgbotapi.BotAPI,
	u tgbotapi.UpdateConfig,
	connPool db.Conn,
) {
	handler := bot.CallbackHandler{botApi, activity.NewService(connPool)}
	msgHandler := bot.NewMessageHandler(botApi, task.NewService(connPool))
	updates := botApi.GetUpdatesChan(u)
	for update := range updates {
		func() {
			ctx, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()

			userID, err := authenticateUser(ctx, update, connPool)
			if err != nil {
				zap.S().Error(err)
			}

			if update.Message == nil {
				if handler.Handle(update) {
					return
				}
				return
			}

			if msgHandler.Handle(&update) {
				return
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			// Old (default) handlers
			switch update.Message.Text {
			case "/start":
				msg.ReplyMarkup = HeadKeyboard
			case CommandHelp:
				s := usecase.NewTaskUseCase(connPool)
				tasks, err := s.GetTasksForUser(ctx, 1)
				if err != nil {
					zap.S().Error(err)
				}
				if len(tasks) == 0 {
					msg.Text = "No tasks?"
					_, err := botApi.Send(msg)
					if err != nil {
						zap.S().Error(err)
					}
				} else {
					for _, t := range tasks {
						tId := strconv.Itoa(t.ID)
						TasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`accept:` + tId)
						TasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`hide:` + tId)
						msg.ReplyMarkup = TasksListKeyboard
						msg.Text = "Task " + tId + "\n" + t.Text
						_, err := botApi.Send(msg)
						if err != nil {
							zap.S().Error(err)
						}
					}
				}
				return
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
			//case CommandGetLocation:
			//	//msg.ReplyMarkup =
			case CommandProcessHelp:
				//msg.ReplyMarkup =
				//msg.Text =
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
				phone, err := getContactsFotUser(ctx, update, userID, connPool)
				if err != nil {
					zap.S().Error(err)
				}
				fmt.Printf("PHONE : %d", phone)
				msg.ReplyMarkup = GetLocationKeyboard
				msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
			}
			if update.Message.Location != nil {
				msg.ReplyMarkup = GetLocationKeyboard
				msg.Text = "Ви поділилися локацією..." + fmt.Sprintf("\nШирота: %v\nДовгота:%v", update.Message.Location.Longitude, update.Message.Location.Latitude)
			}

			_, err = botApi.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}
		}()
	}
}

func authenticateUser(ctx context.Context, update tgbotapi.Update, connPool db.Conn) (int, error) {
	userSocNet := soc_net.UserSocNet{
		UserSocNetID: fmt.Sprintf("%d", update.Message.From.ID),
	}

	s := soc_net.NewService(connPool)

	us, err := s.GetOneBySocNetID(ctx, userSocNet.UserSocNetID)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			su := user.NewService(connPool)
			userID, err := su.CreateOne(ctx, update.Message.From.UserName, 0)
			if err != nil {
				zap.S().Error(err)
			}

			_, err = s.CreateOne(ctx, userID, userSocNet.UserSocNetID)
			if err != nil {
				zap.S().Error(err)
			}
			return us.ID, nil
		}

		zap.S().Error(err)
	}

	return us.ID, err
}

func getContactsFotUser(ctx context.Context, update tgbotapi.Update, userID int, connPool db.Conn) (string, error) {

	su := user.NewService(connPool)
	userID, err := su.UpdateOne(ctx,
		map[string]interface{}{
			`phone_number`: update.Message.Contact.PhoneNumber,
		}, map[string]interface{}{
			`id`: userID,
		})
	if err != nil {
		zap.S().Error(err)

		return "", err
	}

	return update.Message.Contact.PhoneNumber, nil
}
