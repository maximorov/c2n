package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/bootstrap"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/user"
	"helpers/app/domains/user/executor"
	"helpers/app/usecase"
	"regexp"
	"strconv"
	"strings"
)

var coordsRegexp, _ = regexp.Compile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?),\s*[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`)

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
	role            string
}

func (s *MessageHandler) Init() {
	s.handlers = map[string]Handler{
		CommandStart:       &StartHandler{s, HeadKeyboard},
		CommandToMain:      &StartHandler{s, HeadKeyboard},
		CommandInformation: &AboutHandler{s, AfterHeadKeyboard},

		CommandNeedHelp: &NeedHelpHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandCreateNewTask),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		), tgbotapi.NewReplyKeyboard(
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
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		CommandTakeLocationManual: &TakeLocationManualyHandler{s, ToMainKeyboard},
		CommandFiilTaskText:       &WhatFillTaskText{s, ToMainKeyboard},
		CommandMyTasks:            &ShowMyTasksHandler{s, ToMainKeyboard, CancelTaskKeyboard},

		CommandHelp: &HelpHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandTakeNewTask),
				tgbotapi.NewKeyboardButton(CommandMyActiveTasks),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		), tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandTakeNewTask),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		CommandTakeNewTask: &TakeNewTaskHandlerHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation(CommandGetLocation),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		CommandMyActiveTasks: &MyActiveTasksHandler{s, tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandTakeNewTask),
				tgbotapi.NewKeyboardButton(CommandMyActiveTasks),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		SetExecutorLocation: &AfterExecutorLocationSetHandler{s, SetAreaKeyboard},
		CommandRadius1:      &SetRadiusHandler{s, SetAreaKeyboard, 1000},
		CommandRadius3:      &SetRadiusHandler{s, SetAreaKeyboard, 3000},
		CommandRadius5:      &SetRadiusHandler{s, SetAreaKeyboard, 5000},
		CommandRadius10:     &SetRadiusHandler{s, SetAreaKeyboard, 10000},
		CommandNoTasks:      &NoTasksHandler{s, SetAreaKeyboard},
		CommandSubscribe:    &SubscribeHandler{s, AfterHeadKeyboard, true},
		CommandUnsubscribe:  &SubscribeHandler{s, AfterHeadKeyboard, false},
	}
}

func (s *MessageHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	if h, ok := s.handlers[u.Message.Text]; ok {
		h.Handle(ctx, u)
		return
	}

	usr := ctx.Value(`user`).(*user.User)

	switch {
	case u.Message.Contact != nil:
		phone, err := setContactsFotUser(ctx, u, usr.ID, db.GetPool())
		if err != nil {
			zap.S().Error(err)
			// TODO: and what if error?
		}
		if bootstrap.Cnf.Debug {
			zap.S().Debug("PHONE : %d", phone)
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги")
		msg.ReplyMarkup = GetLocationKeyboard
		s.Ans(msg)
	case u.Message.Location != nil:
		switch s.role {
		case "needy":
			err := s.TaskUseCase.CreateRawTask(ctx, usr.ID, u.Message.Location.Latitude, u.Message.Location.Longitude)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, CommandFiilTaskText)
			msg.ParseMode = `markdown`
			msg.ReplyMarkup = ToMainKeyboard
			s.Ans(msg)
		case "executor":
			s.handlers[SetExecutorLocation].Handle(ctx, u)
		default:
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				`Спробуйте спочатку. Геолокацію треба вибирати, як зазначено в інструкції.`,
			)
			msg.ReplyMarkup = HeadKeyboard
			s.Ans(msg)
		}
	default:
		switch { // someone enters coordinates manually
		case coordsRegexp.Match([]byte(u.Message.Text)): // geolocation coordinates
			lonLat := strings.Split(u.Message.Text, `,`)
			lat, _ := strconv.ParseFloat(lonLat[0], 64)
			lon, _ := strconv.ParseFloat(lonLat[1], 64)
			err := s.TaskUseCase.CreateRawTask(ctx, usr.ID, lat, lon)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, CommandFiilTaskText)
			msg.ParseMode = `markdown`
			msg.ReplyMarkup = ToMainKeyboard
			s.Ans(msg)
		default: // any text determines like text of task
			tsk, err := s.TaskService.GetUsersRawTask(ctx, usr.ID)
			if err != nil {
				msg := tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf(`"%s" - Команда не зрозуміла. Спробуйте іншу з варіантів нижче `+SymbLoopDown, u.Message.Text),
				)
				s.Ans(msg)
			} else {
				tuc := usecase.NewTaskUseCase(db.GetPool())
				err = tuc.UpdateLastRawWithText(ctx, tsk.ID, u.Message.Text)
				if err != nil {
					zap.S().Error(err)
				}
				msg := tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf("Ваше завдання #%d\nОчікуйте повідомлення протягом %d годин\nЩойно волонтер візьметься за ваше завдання, ми вам повідомимо.", tsk.ID, task.TaskDeadline),
				)
				msg.ReplyMarkup = ToMainKeyboard
				s.Ans(msg)
			}
		}
	}
}

func (s *MessageHandler) Ans(msg tgbotapi.Chattable) {
	_, err := s.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *MessageHandler) AnsError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, err.Error())
	_, err = s.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
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

func (s *MessageHandler) setLocationFotUser(ctx context.Context, update *tgbotapi.Update, userID int) error {
	point := db.CreatePoint(update.Message.Location.Latitude, update.Message.Location.Longitude)

	_, err := s.ExecutorService.UpdateOne(ctx,
		map[string]interface{}{
			`position`: point,
		}, map[string]interface{}{
			`id`: userID,
		})
	if err != nil {
		zap.S().Error(err)

		return err
	}

	return nil
}
