package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"html"
	"log"
	"net/http"
)

func init() {
	InitEnv("")
	initConfig()
	initLogger()
}

func main() {
	bot, err := tgbotapi.NewBotAPI(Cnf.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go healthcheck()

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			// Extract the command from the Message.
			switch update.Message.Text {
			case "/start":
				msg.ReplyMarkup = HeadKeyboard
			case CommandHelp:
				msg.ReplyMarkup = HelpKeyboard
				msg.Text = "Де вам було б зручно допомогти?"
			case CommandNeedHelp:
				msg.ReplyMarkup = NeedHelpKeyboard
				//msg.Text =
			case CommandInformation:
				msg.Text = Information
			case CommandRadius1:
				msg.ReplyMarkup = GetLocationKeyboard
				//msg.Text =
			case CommandRadius3:
				msg.ReplyMarkup = GetLocationKeyboard
				//msg.Text =
			case CommandRadius5:
				msg.ReplyMarkup = GetLocationKeyboard
				//msg.Text =
			case CommandAllCity:
				//msg.ReplyMarkup =
				//msg.Text =
			case CommandChooseCity:
			//msg.ReplyMarkup =
			//msg.Text =
			case CommandTakeLocationManual:
			//msg.ReplyMarkup =
			//msg.Text =
			case CommandContinueHelp:
				//msg.ReplyMarkup =
				//msg.Text =
			case CommandToMain:
				msg.ReplyMarkup = HeadKeyboard
			case CommandNewTask:
				msg.ReplyMarkup = GetContactsKeyboard
				msg.Text = "Поділіться будь-ласка контактами, щоб з вами могли звʼязатись"
				//TODO: нормально пофиксить эту багу..
			//case CommandGetContact:
			//	msg.ReplyMarkup = GetLocationKeyboard
			//	msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
			case CommandGetLocation:
				//msg.ReplyMarkup =
			case CommandProcessHelp:
				//msg.ReplyMarkup =
				//msg.Text =
			default:

			}

			if update.Message.Contact != nil {
				msg.ReplyMarkup = GetLocationKeyboard
				msg.Text = "Поділіться будь-ласка локацією, щоб з люди знали де ви потребуєте допомоги"
			}
			if update.Message.Location != nil {
				msg.ReplyMarkup = GetLocationKeyboard
				msg.Text = "Ви поділилися локацією..." + fmt.Sprintf("\nШирота: %v\nДовгота:%v", update.Message.Location.Longitude, update.Message.Location.Latitude)
			}

			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func healthcheck() {
	log.Printf("Starts webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
