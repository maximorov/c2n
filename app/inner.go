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
	"helpers/app/domains/task"
	"helpers/app/domains/user"
	"helpers/app/domains/user/executor"
	"helpers/app/domains/user/soc_net"
	"html"
	"log"
	"net/http"
	"strconv"
	"time"
)

func authenticateUser(ctx context.Context, update tgbotapi.Update, sonNetService *soc_net.Service, userService *user.Service, role user.Role) (*user.User, error) {
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

	usr, err := userService.GetOneByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	usr.SetRole(role)

	return usr, nil
}

func informExecutors(ctx context.Context, exRepo *executor.Repository, taskService *task.Service, socNetService *soc_net.Service, callBack *bot.MessageHandler) {
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-t.C:
			n := time.Now()

			if n.Hour() < 8 || n.Hour() > 21 { //at night all sleeps
				continue
			}
			if n.Hour()%4 != 0 || n.Minute() != 0 { //every 3 hours
				continue
			}

			go func(ctx context.Context) {
				executors, err := exRepo.FindMany(ctx,
					[]string{`user_id`, `position`, `area`, `city`, `inform`},
					map[string]interface{}{})
				if err != nil {
					zap.S().Error(err)
				}

				for _, ex := range executors {
					if ex.Inform == false {
						continue
					}

					tasks, err := taskService.FindTasksInRadius(ctx, ex.Position, ex.UserId, float64(ex.Area))
					if err != nil {
						zap.S().Error(err)
					}

					if len(tasks) == 0 {
						continue
					}
					userSocNet, err := socNetService.GetOneByUserID(ctx, ex.UserId)
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
						msg.Text = bot.PrepareTaskText(t)
						callBack.Ans(msg)
					}
					msg.ReplyMarkup = bot.UnsubscribeKeyboard
					msg.Text = `Якщо ви не хочете отримувати автоматичну розсилку, натисніть ` + core.SymbHide
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
				continue
			}

			go func(ctx context.Context) {
				tasks, err := taskRepo.FindMany(ctx,
					[]string{`id`, `user_id`, `position`, `status`, `text`, `deadline`},
					map[string]interface{}{
						`status`: []string{
							task.StatusRaw,
							task.StatusNew,
							task.StatusInProgress,
						},
					})
				if err != nil {
					zap.S().Error(err)
				}

				for _, oneTask := range tasks {
					if oneTask.Deadline.Sub(time.Now()) < 0 {
						_, err = taskRepo.UpdateOne(ctx,
							map[string]interface{}{
								`status`: task.StatusExpired,
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

func healthcheck() {
	log.Printf("Starts webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
