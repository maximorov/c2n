package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/bootstrap"
	"helpers/app/db"
	"html"
	"log"
	"net/http"
)

func init() {
	bootstrap.InitEnv(``)
	bootstrap.InitConfig()
	bootstrap.InitLogger()
}

func main() {
	ctx := context.Background()
	connPool := db.Pool(ctx, bootstrap.Cnf.DB)
	bot, err := tgbotapi.NewBotAPI(bootstrap.Cnf.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go botHandlers(ctx, bot, u, connPool)

	healthcheck()
}

func healthcheck() {
	log.Printf("Starts webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
