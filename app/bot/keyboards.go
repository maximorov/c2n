package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"helpers/app/core"
)

var (
	HeadKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandNeedHelp),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandHelp),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandInformation),
		),
	)
	AfterHeadKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandNeedHelp),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandHelp),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandInformation),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	ToMainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	AttachLocationKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation(CommandGetLocationAuto), // collect location
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandGetLocationHotTo),
		),
		//tgbotapi.NewKeyboardButtonRow(
		//	tgbotapi.NewKeyboardButton(CommandGetLocationManual), // collect location
		//),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	SetAreaKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandRadius1),
			tgbotapi.NewKeyboardButton(CommandRadius3),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandRadius5),
			tgbotapi.NewKeyboardButton(CommandRadius10),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	//GetContactsKeyboard = tgbotapi.NewReplyKeyboard(
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButtonContact(core.CommandGetContact),
	//	),
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton(CommandToMain),
	//	),
	//)
	//GetLocationKeyboard = tgbotapi.NewReplyKeyboard(
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButtonLocation(CommandGetLocationAuto),
	//	),
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton(core.CommandTakeLocationManual),
	//	),
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton(CommandToMain),
	//	),
	//)
	//SubscribeKeyboard = tgbotapi.NewReplyKeyboard(
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton(CommandSubscribe),
	//	),
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButton(CommandToMain),
	//	),
	//)
	UnsubscribeKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandUnsubscribe),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
	TasksListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(core.CommandsWillExecute, ``),
			tgbotapi.NewInlineKeyboardButtonData(core.CommandsRefuseForMe, ``),
		),
	)
	GoogleMapsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(`Google Maps`, `https://www.google.com.ua/maps`),
		),
	)
	ExecutorTasksListKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(core.SymbAccept+` Виконано`, ``),
			tgbotapi.NewInlineKeyboardButtonData(core.SymbRefuse+` Відмовитися`, ``),
		),
	)
	ReopenTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ReopenText, ``),
		),
	)
	CancelTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CancelText, ``),
		),
	)
	SupportInformationKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandSendVideoHowHelp),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandSendVideoHowGetHelp),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandToMain),
		),
	)
)
