package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/user"
	"helpers/app/domains/user/executor"
	"helpers/app/usecase"
	"regexp"
	"strconv"
)

func NewMessageHandler(
	BotApi *tgbotapi.BotAPI,
	ts *task.Service,
	er *executor.Repository,
	es *executor.Service,
	tuc *usecase.TaskUseCase,
) *MessageHandler {
	res := &MessageHandler{
		BotApi:          BotApi,
		TaskService:     ts,
		ExecutorRepo:    er,
		ExecutorService: es,
		TaskUseCase:     tuc,
	}
	res.Init()

	return res
}

type Handler interface {
	Handle(context.Context, *tgbotapi.Update)
}

type MessageHandler struct {
	handlers        map[string]Handler
	BotApi          *tgbotapi.BotAPI
	TaskService     *task.Service
	ExecutorRepo    *executor.Repository
	ExecutorService *executor.Service
	TaskUseCase     *usecase.TaskUseCase
}

func (s *MessageHandler) Init() {
	s.handlers = map[string]Handler{
		CommandStart: &StartHandler{s, HeadKeyboard},

		CommandNeedHelp: &NeedHelpHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandMyTasks),
				tgbotapi.NewKeyboardButton(CommandCreateNewTask),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		CommandCreateNewTask: &CreateTaskHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation(CommandGetLocation), // collect location
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandTakeLocationManual),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		CommandTakeLocationManual: &TaksLocationManualyHandler{s, ToMainKeyboard},
		CommandFiilTaskText:       &WhatFillTaskText{s, ToMainKeyboard},

		SetExecutorLocation: &AfterExecutorLocationSetHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandRadius1),
				tgbotapi.NewKeyboardButton(CommandRadius3),
				tgbotapi.NewKeyboardButton(CommandRadius5),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandAllCity),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		CommandHelp: &HelpHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation(CommandGetLocation), // collect location
			),
			//tgbotapi.NewKeyboardButtonRow(
			//	tgbotapi.NewKeyboardButton(CommandTakeLocationManual),
			//),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		CommandRadius1: &SetRadiusHandler{s, ToMainKeyboard, SetAreaKeyboard, 1000},
		CommandRadius3: &SetRadiusHandler{s, ToMainKeyboard, SetAreaKeyboard, 3000},
		CommandRadius5: &SetRadiusHandler{s, ToMainKeyboard, SetAreaKeyboard, 5000},
		CommandNoTasks: &NoTasksHandler{s, SetAreaKeyboard},
	}
}

func (s *MessageHandler) Handle(ctx context.Context, u *tgbotapi.Update) bool {
	if h, ok := s.handlers[u.Message.Text]; ok {
		h.Handle(ctx, u)
		return true
	}

	//log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	msg.ReplyToMessageID = u.Message.MessageID

	if u.Message.Contact != nil {
		usr := ctx.Value(`user`).(*user.User)
		phone, err := setContactsFotUser(ctx, u, usr.ID, db.GetPool())
		if err != nil {
			zap.S().Error(err)
		}
		fmt.Printf("PHONE : %d", phone)
		msg.ReplyMarkup = GetLocationKeyboard
		msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
	}
	if u.Message.Location != nil {

		msg.Text = ""
		switch u.Message.ReplyToMessage.Text {
		//case CommandCreateNewTask:
		case CommandNeedCollectLocation:
			//err := s.TaskUseCase.CreateRawTask(ctx, usr.ID, u.Message.Location.Longitude, u.Message.Location.Latitude)
			//if err != nil {
			//	zap.S().Error(err)
			//}
			msg.ReplyMarkup = ToMainKeyboard
			msg.Text = CommandFiilTaskText
		case SetExecutorLocation:
			//usr := ctx.Value(`user`).(*user.User)
			s.handlers[SetExecutorLocation].Handle(ctx, u)
			return true
			//err := setLocationFotUser(ctx, u, usr.ID, db.GetPool())
			//if err != nil {
			//	zap.S().Error(err)
			//}
		}
	}

	coordsRegexp, _ := regexp.Compile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?),\s*[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`)
	if coordsRegexp.Match([]byte(u.Message.Text)) { // someone enters coordinates manually
		msg.Text = CommandFiilTaskText
		msg.ReplyMarkup = ToMainKeyboard
		_, err := s.BotApi.Send(msg)
		if err != nil {
			zap.S().Error(err)
		}
		return true
	}

	switch u.Message.Text {
	case CommandInformation:
		msg.Text = Information
	case CommandAllCity:
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
	default: // any text determines like text of task
		//usr := ctx.Value(`user`).(*user.User)
		//tsk, err := s.TaskService.GetUsersRawTask(ctx, usr.ID)
		//if err != nil {
		//	msg.Text += "Your RAW task didnt found. Start from the beginning"
		//	_, err = s.BotApi.Send(msg)
		//	if err != nil {
		//		zap.S().Error(err)
		//	}
		//	return true
		//}
		//s := usecase.NewTaskUseCase(db.GetPool())
		//err = s.UpdateLastRawWithText(ctx, tsk.ID, u.Message.Text)
		//if err != nil {
		//	zap.S().Error(err)
		//}
		//msg.Text = "Your task #" + strconv.Itoa(tsk.ID) + "\n" +
		msg.Text = "Your task #111" + "\n" +
			" свами должны связаться в течении " + strconv.Itoa(task.TaskDeadline) + " часов"
		msg.ReplyMarkup = ToMainKeyboard
	}

	_, err := s.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}

	return false
}

func setContactsFotUser(ctx context.Context, update *tgbotapi.Update, userID int, connPool db.Conn) (string, error) {
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

func setLocationFotUser(ctx context.Context, update *tgbotapi.Update, userID int, connPool db.Conn) error {
	s := executor.NewService(connPool)
	ex, err := s.GetOneByUserID(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	point := db.CreatePoint(update.Message.Location.Latitude, update.Message.Location.Longitude)

	_, err = s.UpdateOne(ctx,
		map[string]interface{}{
			`position`: point,
		}, map[string]interface{}{
			`id`: ex.ID,
		})
	if err != nil {
		zap.S().Error(err)

		return err
	}

	return nil
}
