package main

import (
	"context"
	"errors"
	"fmt"
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
	"time"
)

func botHandlers(
	ctx context.Context,
	botApi *tgbotapi.BotAPI,
	u tgbotapi.UpdateConfig,
	connPool db.Conn,
) {
	taskServie := task.NewService(connPool)

	clbkHandler := bot.CallbackHandler{
		botApi,
		activity.NewService(connPool),
	}
	msgHandler := bot.NewMessageHandler(
		botApi,
		taskServie,
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
