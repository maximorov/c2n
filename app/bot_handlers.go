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
	"strconv"
	"time"
)

func botHandlers(
	ctx context.Context,
	botApi *tgbotapi.BotAPI,
	u tgbotapi.UpdateConfig,
	connPool db.Conn,
) {
	taskServie := task.NewService(connPool)

	handler := bot.CallbackHandler{botApi, activity.NewService(connPool)}
	msgHandler := bot.NewMessageHandler(botApi, taskServie)
	updates := botApi.GetUpdatesChan(u)
	for update := range updates {
		func() {
			ctx, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()

			usr, err := authenticateUser(ctx, update, connPool)
			if err != nil {
				zap.S().Error(err)
			}

			ctx = context.WithValue(ctx, `user`, usr)

			if update.Message == nil {
				if handler.Handle(update) {
					return
				}
				return
			}

			if msgHandler.Handle(ctx, &update) {
				return
			}

			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			// Old (default) handlers
			switch update.Message.Text {
			case "/start":
				msg.ReplyMarkup = bot.HeadKeyboard
			case bot.CommandHelp:
				//er := executor.NewRepo(connPool)
				//s := usecase.NewTaskUseCase(connPool)
				//executor, err := er.FindOne(ctx, []string{`position`}, map[string]interface{}{`user_id`: usr.ID})
				//tasks, err := s.GetTasksForUser(ctx, executor.Position)
				//if err != nil {
				//	zap.S().Error(err)
				//}
				//if len(tasks) == 0 {
				//	msg.Text = "No tasks?"
				//	_, err := botApi.Send(msg)
				//	if err != nil {
				//		zap.S().Error(err)
				//	}
				//} else {
				//	for _, t := range tasks {
				//		tId := strconv.Itoa(t.ID)
				//		bot.TasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`accept:` + tId)
				//		bot.TasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`hide:` + tId)
				//		msg.ReplyMarkup = bot.TasksListKeyboard
				//		msg.Text = "Task " + tId + "\n" + t.Text
				//		_, err := botApi.Send(msg)
				//		if err != nil {
				//			zap.S().Error(err)
				//		}
				//	}
				//}
				return
			case bot.CommandInformation:
				msg.Text = bot.Information
			case bot.CommandRadius1:
				msg.ReplyMarkup = bot.GetLocationKeyboard
				//msg.Text =
			case bot.CommandRadius3:
				msg.ReplyMarkup = bot.GetLocationKeyboard
				//msg.Text =
			case bot.CommandRadius5:
				msg.ReplyMarkup = bot.GetLocationKeyboard
				//msg.Text =
			case bot.CommandAllCity:
				//msg.ReplyMarkup =
				//msg.Text =
			case bot.CommandChooseCity:
			//msg.ReplyMarkup =
			//msg.Text =
			case bot.CommandTakeLocationManual:
			//msg.ReplyMarkup =
			//msg.Text =
			case bot.CommandContinueHelp:
				//msg.ReplyMarkup =
				//msg.Text =
			case bot.CommandToMain:
				msg.ReplyMarkup = bot.HeadKeyboard
			case bot.CommandNewTask:
				msg.ReplyMarkup = bot.GetContactsKeyboard
				msg.Text = "Поділіться будь-ласка контактами, щоб з вами могли звʼязатись"
				//TODO: нормально пофиксить эту багу..
			//case CommandGetContact:
			//	msg.ReplyMarkup = GetLocationKeyboard
			//	msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
			//case CommandGetLocation:
			//	//msg.ReplyMarkup =
			case bot.CommandProcessHelp:
				//msg.ReplyMarkup =
				//msg.Text =
			case bot.CommandCreateTask:
				s := task.NewService(connPool)
				uId := 1
				x, y := 1.0, 1.0
				taskId, err := s.CreateTask(ctx, uId, x, y, update.Message.Text)
				if err != nil {
					zap.S().Error(err)
				}
				//msg.ReplyMarkup =
				msg.Text = strconv.Itoa(taskId)
			default: // any text determines like text of task
				tsk, err := taskServie.GetUsersRawTask(ctx, usr.ID)
				if err != nil {
					msg.Text += "Your RAW task didnt found. Start from the beginning"
					_, err = botApi.Send(msg)
					if err != nil {
						zap.S().Error(err)
					}
					return
				}
				s := usecase.NewTaskUseCase(connPool)
				err = s.UpdateLastRawWithText(ctx, tsk.ID, update.Message.Text)
				if err != nil {
					zap.S().Error(err)
				}
				msg.Text = "Your task #" + strconv.Itoa(tsk.ID) + "\n" +
					" свами должны связаться в течении " + strconv.Itoa(task.TaskDeadline) + " часов"
			}

			if update.Message.Contact != nil {
				phone, err := getContactsFotUser(ctx, update, usr.ID, connPool)
				if err != nil {
					zap.S().Error(err)
				}
				fmt.Printf("PHONE : %d", phone)
				msg.ReplyMarkup = bot.GetLocationKeyboard
				msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
			}
			if update.Message.Location != nil {
				msg.Text = ""
				//msg.Text = "Ви поділилися локацією..." + fmt.Sprintf("\nШирота: %v\nДовгота:%v", update.Message.Location.Longitude, update.Message.Location.Latitude)
				switch update.Message.ReplyToMessage.Text {
				case `We need to collect info about you`:
					s := usecase.NewTaskUseCase(connPool)
					err := s.CreateRawTask(ctx, usr.ID, update.Message.Location.Longitude, update.Message.Location.Latitude)
					if err != nil {
						zap.S().Error(err)
					}
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					msg.Text += "\nWrite your problem"
				}
			}

			_, err = botApi.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}
		}()
	}
}

func authenticateUser(ctx context.Context, update tgbotapi.Update, connPool db.Conn) (*user.User, error) {
	if update.Message == nil {
		return nil, nil
	}

	userSocNet := soc_net.UserSocNet{
		UserSocNetID: fmt.Sprintf("%d", update.Message.From.ID),
	}

	var userId int

	s := soc_net.NewService(connPool)
	su := user.NewService(connPool)

	us, err := s.GetOneBySocNetID(ctx, userSocNet.UserSocNetID)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			userID, err := su.CreateOne(ctx, update.Message.From.UserName, "")
			if err != nil {
				return nil, err
			}

			userId, err = s.CreateOne(ctx, userID, userSocNet.UserSocNetID)
			if err != nil {
				return nil, err
			}
		}

		return nil, err
	} else {
		userId = us.UserId
	}

	u, err := su.GetOneByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return u, nil
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

func getLocationFotUser(ctx context.Context, update tgbotapi.Update, userID int, connPool db.Conn) (string, error) {

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
