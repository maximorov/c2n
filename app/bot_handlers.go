package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"helpers/app/bot"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/task/activity"
	"helpers/app/domains/user"
	"helpers/app/domains/user/executor"
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
	clbkHandler := bot.CallbackHandler{
		botApi,
		activity.NewService(connPool),
	}
	msgHandler := bot.NewMessageHandler(
		botApi,
		task.NewService(connPool),
		executor.NewRepo(connPool),
		executor.NewService(connPool),
		usecase.NewTaskUseCase(connPool),
	)

	updates := botApi.GetUpdatesChan(u)
	for update := range updates {
		func() {
			ctx, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()

			usr, err := authenticateUser(ctx, update, connPool)
			if err != nil {
				zap.S().Error(err)
				msgHandler.AnsError(update.Message.Chat.ID, err)
				return
			}

			ctx = context.WithValue(ctx, `user`, usr)

			switch {
			case update.Message != nil:
				msgHandler.Handle(ctx, &update)
			case update.CallbackData() != ``:
				clbkHandler.Handle(update)
			}
		}()
	}
}

func authenticateUser(ctx context.Context, update tgbotapi.Update, connPool db.Conn) (*user.User, error) {
	if update.Message == nil {
		return nil, nil // TODO: process when no message
	}

	usrSocId := strconv.Itoa(int(update.Message.From.ID))
	var userId int

	sonNetService := soc_net.NewService(connPool)
	userService := user.NewService(connPool)

	usrSoc, err := sonNetService.GetOneBySocNetID(ctx, usrSocId)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			userId, err = userService.CreateOne(ctx, update.Message.From.UserName, "")
			if err != nil {
				return nil, err
			}

			_, err = sonNetService.CreateOne(ctx, userId, usrSocId)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		userId = usrSoc.UserId
	}

	u, err := userService.GetOneByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return u, nil
}
