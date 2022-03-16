package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/domains/task/activity"
	"helpers/app/domains/user"
	"helpers/app/usecase"
	"strconv"
	"strings"
)

const ReopenText = `Reopen`
const ReopenCallback = `reopen`

type CallbackHandler struct {
	BotApi              *tgbotapi.BotAPI
	TaskActivityService *activity.Service
	TaskUseCase         *usecase.TaskUseCase
}

func (s *CallbackHandler) Ans(msg tgbotapi.MessageConfig) {
	_, err := s.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *CallbackHandler) AnsError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, err.Error())
	_, err = s.BotApi.Send(msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (s *CallbackHandler) Handle(ctx context.Context, u tgbotapi.Update) bool {
	usr := ctx.Value(`user`).(*user.User)
	handled := false

	switch {
	case strings.Contains(u.CallbackData(), `hide:`) || strings.Contains(u.CallbackData(), `accept:`):
		handled = true
		parsed := strings.Split(u.CallbackData(), `:`)
		action := parsed[0]
		taskId, _ := strconv.Atoi(parsed[1])

		switch action {
		case `accept`:
			err := s.TaskActivityService.CreateActivity(ctx, usr.ID, taskId, `taken`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Завдання треба виконати за добу`)
			s.Ans(msg)
			s.informNeedy(ctx, taskId, `taken`)
		case `hide`:
			err := s.TaskActivityService.CreateActivity(ctx, usr.ID, taskId, `hidden`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Це завдання візьме хтось інший`)
			s.Ans(msg)
		}
	case strings.Contains(u.CallbackData(), `complete:`) || strings.Contains(u.CallbackData(), `refuse:`):
		handled = true
		parsed := strings.Split(u.CallbackData(), `:`)
		action := parsed[0]
		taskId, _ := strconv.Atoi(parsed[1])

		switch action {
		case `complete`:
			err := s.TaskActivityService.UpdateActivity(ctx, usr.ID, taskId, `completed`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Дякую`)
			s.Ans(msg)
			s.informNeedy(ctx, taskId, `complete`)
		case `refuse`:
			err := s.TaskActivityService.UpdateActivity(ctx, usr.ID, taskId, `refused`)
			if err != nil {
				zap.S().Error(err)
			}
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Візьме хтось інший, чи ні`)
			s.Ans(msg)
			s.informNeedy(ctx, taskId, `refuse`)
		}
	}

	return handled
}

func (s *CallbackHandler) informNeedy(ctx context.Context, tId int, status string) {
	socUser, err := s.TaskUseCase.GetSocUserByTask(ctx, tId)
	if err != nil {
		zap.S().Error(err)
		return
	}
	chatId, err := strconv.Atoi(socUser.SocNetID)
	if err != nil {
		zap.S().Error(err)
		return
	}

	kb := ReopenTaskKeyboard
	kb.InlineKeyboard[0][0].CallbackData = core.StrP(ReopenCallback + `:` + strconv.Itoa(tId))

	msg := tgbotapi.NewMessage(int64(chatId), ``)
	msg.ReplyMarkup = kb

	switch status {
	case `refuse`:
		msg.Text = `Волонтер отказался от вашей задачи. Ждем другого`
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, `new`)
		if err != nil {
			zap.S().Error(err)
			return
		}
	case `complete`:
		msg.Text = `Волонтер отметил вашу задачу как выполнено. Если это не так - нажмите ` + ReopenText
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, `done`)
		if err != nil {
			zap.S().Error(err)
			return
		}
	case `taken`:
		msg.Text = `Волонтер вляз вашу задачу в работу. Если в течении нескольких часов с вами не связались, нажмите ` + ReopenText
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, `in_progress`)
		if err != nil {
			zap.S().Error(err)
			return
		}
	}

	s.Ans(msg)
}
