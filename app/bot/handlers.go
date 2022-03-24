package bot

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/task/activity"
	"helpers/app/domains/user"
	"helpers/app/domains/user/executor"
	"helpers/app/domains/user/soc_net"
	"helpers/app/usecase"
	"regexp"
	"strconv"
	"time"
)

var coordsRegexp, _ = regexp.Compile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?),\s*[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`)

func NewMessageHandler(
	BotApi *tgbotapi.BotAPI,
	ts *task.Service,
	tas *activity.Service,
	er *executor.Repository,
	es *executor.Service,
	tuc *usecase.TaskUseCase,
	snr *soc_net.Repository,
) *MessageHandler {
	res := &MessageHandler{
		BotApi:              BotApi,
		TaskService:         ts,
		TaskActivityService: tas,
		ExecutorRepo:        er,
		ExecutorService:     es,
		TaskUseCase:         tuc,
		SocNetRepoRepo:      snr,
	}
	res.Init()

	return res
}

type Handler interface {
	Handle(context.Context, *tgbotapi.Update) bool
	UserRole() user.Role
}

type MessageHandler struct {
	handlers            map[string]Handler
	replyHandlers       map[string]Handler
	mainTextHandler     Handler
	callbackHandler     Handler
	BotApi              *tgbotapi.BotAPI
	TaskService         *task.Service
	TaskActivityService *activity.Service
	ExecutorRepo        *executor.Repository
	SocNetRepoRepo      *soc_net.Repository
	ExecutorService     *executor.Service
	TaskUseCase         *usecase.TaskUseCase
}

func (s *MessageHandler) Init() {
	s.handlers = map[string]Handler{
		CommandStart:       &StartHandler{s, HeadKeyboard},
		CommandToMain:      &StartHandler{s, HeadKeyboard},
		CommandInformation: &AboutHandler{s, SupportInformationKeyboard},

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
				tgbotapi.NewKeyboardButtonLocation(CommandGetLocationAuto), // collect location
			),
			//tgbotapi.NewKeyboardButtonRow(
			//	tgbotapi.NewKeyboardButton(CommandGetLocationManual), // collect location
			//),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(CommandToMain),
			),
		)},
		//CommandGetLocationManual: &TakeLocationManualyHandler{s, ToMainKeyboard},
		CommandFiilTaskText: &WhatFillTaskText{s, ToMainKeyboard},
		CommandMyTasks:      &ShowMyTasksHandler{s, ToMainKeyboard, CancelTaskKeyboard},

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
				tgbotapi.NewKeyboardButtonLocation(CommandGetLocationAuto),
			),
			//tgbotapi.NewKeyboardButtonRow(
			//	tgbotapi.NewKeyboardButton(CommandGetLocationManual), // collect location
			//),
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
		CommandRadius1:  &SetRadiusHandler{s, SetAreaKeyboard, 1000},
		CommandRadius3:  &SetRadiusHandler{s, SetAreaKeyboard, 3000},
		CommandRadius5:  &SetRadiusHandler{s, SetAreaKeyboard, 5000},
		CommandRadius10: &SetRadiusHandler{s, SetAreaKeyboard, 10000},
		//CommandNoTasks:      &NoTasksHandler{s, SetAreaKeyboard},
		CommandSubscribe:           &SubscribeHandler{s, AfterHeadKeyboard, true},
		CommandUnsubscribe:         &SubscribeHandler{s, AfterHeadKeyboard, false},
		CommandSendVideoHowHelp:    &SupportInformationHendler{s, SupportInformationKeyboard, CommandSendVideoHowHelp},
		CommandSendVideoHowGetHelp: &SupportInformationHendler{s, SupportInformationKeyboard, CommandSendVideoHowGetHelp},
	}

	s.replyHandlers = map[string]Handler{
		DoNotGiveLocationNeedy:    &SetTaskLocationHandler{s, ToMainKeyboard},
		DoNotGiveLocationExecutor: &SetExecutorLocationHandler{s, SetAreaKeyboard},
	}

	s.mainTextHandler = &MainTextHandler{s, ToMainKeyboard}

	s.callbackHandler = &CallbackHandler{s}
}

func (s *MessageHandler) GetUserRole(ctx context.Context, u *tgbotapi.Update) user.Role {
	h := s.DetectHandler(ctx, u)
	if h != nil {
		return h.UserRole()
	}

	return user.Unknown
}

func (s *MessageHandler) DetectHandler(ctx context.Context, u *tgbotapi.Update) Handler {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Error("Recovering from panic:", r)
		}
	}()

	if u.CallbackData() != `` {
		return s.callbackHandler
	}

	if u.Message == nil {
		zap.S().Error(`Someone delete bot`)
		return nil
	}

	if u.Message.ReplyToMessage != nil {
		if h, ok := s.replyHandlers[u.Message.ReplyToMessage.Text]; ok {
			return h
		}
	}

	snu, err := s.SocNetRepoRepo.FindOne(
		ctx,
		[]string{`last_received_message`},
		sq.Eq{`soc_net_id`: strconv.Itoa(int(u.Message.Chat.ID))},
	)
	if err != nil && !errors.As(err, &pgx.ErrNoRows) {
		zap.S().Error(err)
	}
	if snu != nil && snu.LastReceivedMessage != `` {
		if h, ok := s.replyHandlers[snu.LastReceivedMessage]; ok {
			return h
		}
	}

	if h, ok := s.handlers[u.Message.Text]; ok {
		return h
	}

	return s.mainTextHandler
}

func (s *MessageHandler) Handle(ctx context.Context, u *tgbotapi.Update) {
	handler := s.DetectHandler(ctx, u)
	if handler != nil && !handler.Handle(ctx, u) {
		// Else events handler
		msg := tgbotapi.NewMessage(
			u.Message.Chat.ID,
			`Команда не зрозуміла. Виберіть одну з тих, що нижче. `+SymbLoopDown,
		)
		//msg.ReplyMarkup = HeadKeyboard
		s.Ans(msg)
	}
}

func (s *MessageHandler) Ans(msg tgbotapi.Chattable) {
	_, err := s.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}

	if mc, ok := msg.(tgbotapi.MessageConfig); ok {
		if mc.Text != `` {
			_, err = s.SocNetRepoRepo.UpdateOne(
				context.Background(),
				map[string]interface{}{`last_received_message`: mc.Text},
				sq.Eq{`soc_net_id`: strconv.Itoa(int(mc.ChatID))},
			)
			if err != nil {
				zap.S().Error(err)
			}
		}
	}

}

func (s *MessageHandler) AnsError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, err.Error())
	_, err = s.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *MessageHandler) AnsDelete(msg tgbotapi.DeleteMessageConfig) {
	_, err := s.BotApi.Request(msg)
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

func PrepareTaskText(task *task.Task) string {
	past := time.Since(task.Created)
	hoursAgo := past.Hours()
	var pastText string
	if hoursAgo < 1 {
		pastText = `менш ніж годину тому`
	} else {
		pastText = strconv.Itoa(int(hoursAgo)) + ` годин тому`
	}

	result := fmt.Sprintf(
		"%s Завдання #%d\nСтворено %s\n\n%s\n\nВідстань від вас: %.0f метрів",
		SymbTask, task.ID, pastText, task.Text, task.GetDistance())

	return result
}
