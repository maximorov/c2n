package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	HeadKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandHelp),
			tgbotapi.NewKeyboardButton(CommandNeedHelp),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandInformation),
		),
	)
	ToMainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	SetAreaKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandRadius1),
			tgbotapi.NewKeyboardButton(CommandRadius3),
			tgbotapi.NewKeyboardButton(CommandRadius5),
		),
		//tgbotapi.NewKeyboardButtonRow(
		//	tgbotapi.NewKeyboardButton(CommandAllCity),
		//),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	GetContactsKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(CommandGetContact),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	GetLocationKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation(CommandGetLocation),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandTakeLocationManual),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	TasksListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CommandsWillExecute, ``),
			tgbotapi.NewInlineKeyboardButtonData(CommandsRefuseForMe, ``),
		),
	)
	ExecutorTasksListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(`Виконав`, ``),
			tgbotapi.NewInlineKeyboardButtonData(`Відмовляюсь`, ``),
		),
	)
	ReopenTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ReopenText, ``),
		),
	)
)
