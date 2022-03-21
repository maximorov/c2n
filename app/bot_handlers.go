package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"helpers/app/bot"
	"helpers/app/core"
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
	exRepo := executor.NewRepo(connPool)
	taskRepo := task.NewRepo(connPool)
	tUC := usecase.NewTaskUseCase(connPool)

	clbkHandler := bot.CallbackHandler{
		botApi,
		activity.NewService(connPool),
		tUC,
	}
	msgHandler := bot.NewMessageHandler(
		botApi,
		task.NewService(connPool),
		exRepo,
		executor.NewService(connPool),
		tUC,
	)

	go informExecutors(ctx, exRepo, connPool, clbkHandler)
	go setExpired(ctx, taskRepo)

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
				clbkHandler.Handle(ctx, update)
			}
		}()
	}
}

func authenticateUser(ctx context.Context, update tgbotapi.Update, connPool db.Conn) (*user.User, error) {
	var fromId int

	switch {
	case update.Message != nil:
		fromId = int(update.Message.From.ID)
	case update.CallbackQuery != nil:
		fromId = int(update.CallbackQuery.From.ID)
	default:
		zap.S().Error(`Ahtung! User not authorized`)
		return nil, nil
	}

	usrSocId := strconv.Itoa(fromId)
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

func informExecutors(ctx context.Context, exRepo *executor.Repository, connPool db.Conn, callBack bot.CallbackHandler) {
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-t.C:
			n := time.Now()

			if n.Hour() < 8 || n.Hour() > 22 { //at night all sleeps
				continue
			}
			if n.Hour()%3 != 0 || n.Minute() != 0 { //every 3 hours
				//zap.S().Info("Not a time for inform executors")
				continue
			}

			//zap.S().Info("Time for inform executors!")

			go func(ctx context.Context) {
				executors, err := exRepo.FindMany(ctx,
					[]string{`user_id`, `position`, `area`, `city`, `inform`},
					map[string]interface{}{})
				if err != nil {
					zap.S().Error(err)
				}

				for _, executor := range executors {
					if executor.Inform == false {
						continue
					}

					s := task.NewService(connPool)
					tasks, err := s.FindTasksInRadius(ctx, executor.Position, executor.UserId, float64(executor.Area))
					if err != nil {
						zap.S().Error(err)
					}

					if len(tasks) == 0 {
						continue
					}
					sSoc := soc_net.NewService(connPool)
					userSocNet, err := sSoc.GetOneByUserID(ctx, executor.UserId)
					if err != nil {
						zap.S().Error(err)
					}

					socNenID, _ := strconv.Atoi(userSocNet.SocNetID)

					msg := tgbotapi.NewMessage(int64(socNenID), "Люди на обраній вами території потребують допомоги:")
					callBack.Ans(msg)

					for _, t := range tasks {
						tId := strconv.Itoa(t.ID)
						bot.TasksListKeyboard.InlineKeyboard[0][0].CallbackData = core.StrP(`accept:` + tId)
						bot.TasksListKeyboard.InlineKeyboard[0][1].CallbackData = core.StrP(`hide:` + tId)
						msg.ReplyMarkup = bot.TasksListKeyboard
						msg.Text = "Task " + tId + "\n" + t.Text
						callBack.Ans(msg)
					}
					msg.ReplyMarkup = bot.UnsubscribeKeyboard
					msg.Text = "Якщо не хочете отримувати автоматичну розсилку натисніть \"Відписатися\""
					callBack.Ans(msg)
				}
			}(ctx)
		}
	}
}

func setExpired(ctx context.Context, taskRepo *task.Repository) {
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-t.C:
			n := time.Now()

			if n.Minute() != 0 { //every hour
				zap.S().Info("Not a time for expire tasks")
				continue
			}

			zap.S().Info("Time for expire tasks!")

			go func(ctx context.Context) {
				tasks, err := taskRepo.FindMany(ctx,
					[]string{`id`, `user_id`, `position`, `status`, `text`, `deadline`},
					map[string]interface{}{})
				if err != nil {
					zap.S().Error(err)
				}

				for _, oneTask := range tasks {
					if oneTask.Deadline.Sub(time.Now()) < 0 {
						_, err = taskRepo.UpdateOne(ctx,
							map[string]interface{}{
								`status`: "expired",
							}, map[string]interface{}{
								`id`: oneTask.ID,
							})
						if err != nil {
							zap.S().Error(err)
						}
					}
				}
			}(ctx)
		}
	}
}
