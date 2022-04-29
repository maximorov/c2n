package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"helpers/app/bootstrap"
	"helpers/app/bot"
	"helpers/app/core/db"
	"helpers/app/domains/task"
	"helpers/app/domains/task/activity"
	"helpers/app/domains/user"
	"helpers/app/domains/user/executor"
	"helpers/app/domains/user/soc_net"
	"helpers/app/usecase"
	"log"
	"time"
)

func init() {
	bootstrap.InitEnv(``)
	bootstrap.InitConfig()
	bootstrap.InitLogger()
}

func main() {
	ctx := context.Background()
	botApi, err := tgbotapi.NewBotAPI(bootstrap.Cnf.TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	botApi.Debug = bootstrap.Cnf.Debug
	if bootstrap.Cnf.Debug {
		log.Printf("Authorized on account %s", botApi.Self.UserName)
	}

	connPool := db.Pool(ctx, bootstrap.Cnf.DB)
	exService := executor.NewService(connPool)
	taskService := task.NewService(connPool)
	tUC := usecase.NewTaskUseCase(connPool)
	socNetService := soc_net.NewService(connPool)
	sonNetService := soc_net.NewService(connPool)
	userService := user.NewService(connPool)
	activiryService := activity.NewService(connPool)

	msgHandler := bot.NewMessageHandler(
		botApi,
		taskService,
		activiryService,
		exService,
		tUC,
		socNetService.Repo,
	)

	go informExecutors(ctx, exService.Repo, taskService, socNetService, msgHandler)
	go setExpired(ctx, taskService.Repo)
	go healthcheck()

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60
	updates := botApi.GetUpdatesChan(update)
	for u := range updates {
		func(u tgbotapi.Update) {
			ctx, cancel := context.WithTimeout(ctx, time.Second*30)
			defer cancel()

			usr, err := authenticateUser(ctx, u, sonNetService, userService, msgHandler.GetUserRole(ctx, &u))
			if err != nil {
				zap.S().Error(err)
				msgHandler.AnsError(u.Message.Chat.ID, err)
				return
			}

			ctx = context.WithValue(ctx, `user`, usr)

			msgHandler.Handle(ctx, &u)
		}(u)
	}
}
