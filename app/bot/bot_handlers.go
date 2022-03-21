package bot

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"helpers/app/core"
	"helpers/app/domains/task"
	"helpers/app/domains/task/activity"
	"helpers/app/domains/user"
	"helpers/app/usecase"
	"strconv"
	"strings"
	"time"
)

const ReopenText = SymbClapper + ` Перевідкрити`
const ReopenCallback = `reopen`

const CancelText = SymbRefuse + ` Видалити завдання`
const CancelCallback = `cancel`

type CallbackHandler struct {
	BotApi              *tgbotapi.BotAPI
	TaskActivityService *activity.Service
	TaskUseCase         *usecase.TaskUseCase
}

func (s *CallbackHandler) AnsDelete(msg tgbotapi.DeleteMessageConfig) {
	_, err := s.BotApi.Request(msg)
	if err != nil {
		zap.S().Error(err)
	}
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

	callbackData := u.CallbackData()
	if !strings.Contains(callbackData, `:`) {
		return false
	}

	handled = true
	parsed := strings.Split(u.CallbackData(), `:`)
	action := parsed[0]
	taskId, _ := strconv.Atoi(parsed[1])

	t, err := s.TaskUseCase.TaskRepo.FindOne(ctx, []string{`deadline`}, sq.Eq{`id`: taskId})
	if err != nil {
		zap.S().Error(err)
		return false
	}
	timePast := t.Deadline.Sub(time.Now()).Hours()

	switch action {
	case `accept`:
		if _, err := s.TaskActivityService.Repo.FindOne(ctx, []string{`task_id`}, sq.Eq{`task_id`: taskId, `executor_id`: usr.ID}); err != nil {
			err := s.TaskActivityService.CreateActivity(ctx, usr.ID, taskId, `taken`)
			if err != nil {
				zap.S().Error(err)
			}
		} else {
			err := s.TaskActivityService.UpdateActivity(ctx, usr.ID, taskId, `taken`)
			if err != nil {
				zap.S().Error(err)
			}
		}

		msgDel := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
		s.AnsDelete(msgDel)
		msg := tgbotapi.NewMessage(
			u.CallbackQuery.Message.Chat.ID,
			fmt.Sprintf("%s Завдання #%d потрібно виконати за %d годин. На вашу допомогу вже чекають.\n", SymbHart, taskId, timePast),
		)
		s.Ans(msg)
		s.informNeedy(ctx, taskId, `taken`)
	case `hide`:
		err := s.TaskActivityService.CreateActivity(ctx, usr.ID, taskId, `hidden`)
		if err != nil {
			zap.S().Error(err)
			return true
		}
		msgDel := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
		s.AnsDelete(msgDel)
		msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Це завдання візьме інший волонтер.`)
		s.Ans(msg)
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
		msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Ви відмовилися від виконання цього завдання, відтепер воно доступне для іншого волонтера.`)
		s.Ans(msg)
		s.informNeedy(ctx, taskId, `refuse`)
	case `reopen`:
		s.informNeedy(ctx, taskId, `reopen`)
	case CancelCallback:
		if s.informExecutor(ctx, taskId, CancelCallback) {
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Виконався попереджено`)
			s.Ans(msg)
		} else {
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, `Завдання видалено.`)
			s.Ans(msg)
		}
	default:
		zap.S().Error(`callback didnt processed`)
		handled = false
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
		msg.Text = `На жаль, цей волонтер відмовився від виконання цього завдання. Триває пошук іншого волонтера.`
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, task.StatusNew)
		if err != nil {
			zap.S().Error(err)
			return
		}
	case `complete`:
		msg.Text = fmt.Sprintf(SymbTask+" Завдання #%d\n"+
			"Позначено виконаним. Якщо це не так, натисніть на кнопку\n"+
			"[%s]", tId, ReopenText)
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, task.StatusDone)
		if err != nil {
			zap.S().Error(err)
			return
		}
	case `taken`:
		msg.Text = fmt.Sprintf(SymbTask+" Завдання %d\n"+
			"Хтось узяв його у роботу.\n"+
			"Якщо протягом декількох годин ви не отримали повідомлення від волонтера, "+
			"натисніть на кнопку %s - ми будемо шукати іншого.", tId, ReopenText)
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, task.StatusInProgress)
		if err != nil {
			zap.S().Error(err)
			return
		}
	case `reopen`:
		msg.Text = fmt.Sprintf("Ваше завдання #%d чекає на іншого волонтера", tId)
		msg.ReplyMarkup = ToMainKeyboard
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, task.StatusNew)
		if err != nil {
			zap.S().Error(err)
			return
		}
	}

	s.Ans(msg)
}

// returns is executor wos informed
func (s *CallbackHandler) informExecutor(ctx context.Context, tId int, status string) bool {
	switch status {
	case CancelCallback:
		err := s.TaskUseCase.UpdateTaskStatus(ctx, tId, task.StatusCancelled)
		if err != nil {
			zap.S().Error(err)
			return false
		}
	}

	socExecutor, err := s.TaskUseCase.GetSocExecutorByTaskActivity(ctx, tId)
	if socExecutor == nil {
		if err != nil && !errors.As(err, &pgx.ErrNoRows) {
			zap.S().Error(err)
		}
		return false
	}
	chatId, err := strconv.Atoi(socExecutor.SocNetID)
	if err != nil {
		zap.S().Error(err)
		return false
	}

	msg := tgbotapi.NewMessage(int64(chatId), ``)
	msg.ReplyMarkup = ToMainKeyboard
	msg.Text = `Задача закрыта. Составитель отказался`

	s.Ans(msg)

	return true
}
