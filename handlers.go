package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//writes pong to chat
func pingHandler(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	resp := tgbotapi.NewMessage(upd.Message.Chat.ID, "pong")

	_, e := bot.Send(resp)

	if e != nil {
		log.Panic(e)
	}
}

//add user to database and create user settings
//'/start' - command
func startHandler(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) {

}

func listHandler(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) {

}

//show available commands and usage examples
//'/help' - command
func helpHandler(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	str := `available commands
	/remind 02.01.2006_15:04 [often] [text] - add reminder
	/settings - list or chage settings for current user
	/start - add user to DB. After this operation settings can be changed
	/help - print this message
	/ping - bot responds with 'pong'`

	bot.Send(tgbotapi.NewMessage(upd.Message.Chat.ID, str))
}

//for test of current feature
//'/test' - command
func testHandler(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) {

	var arr []tgbotapi.KeyboardButton

	arr = append(arr, tgbotapi.NewKeyboardButton("hello"))

	kb := tgbotapi.NewReplyKeyboard(arr)

	kb.OneTimeKeyboard = true

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "a")
	msg.ReplyMarkup = kb

	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	_, e := bot.Send(msg)

	if e != nil {
		log.Println("failed to send kb: ", e)
	}
}

//show user settings and allows to change them
//'/settings' - command
func settingsHandler(msg *tgbotapi.Update, bot *tgbotapi.BotAPI) {

}

//main functionality
//'/remind' - command
//'/remind [when] [how often] [what]
//how often may be one of the following
//once - reminder will be deleted after remind
//daily
//example: /remind 11.12.2020_2:00 once suck a dick
func setReminderHandler(msg *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	settings, err := mongoCtl.findUserSettings(msg.Message.Chat.UserName)

	if err != nil {
		settings = getDefaultUserSettings()
	}

	t := extractFirst(&msg.Message.Text)
	i := extractFirst(&msg.Message.Text)

	l, err := time.LoadLocation(settings.Location)

	if err != nil {
		l, _ = time.LoadLocation("Europe/Moscow")
	}

	date, err := time.ParseInLocation(settings.Layout, t, l)

	if err != nil {
		bot.Send(tgbotapi.NewMessage(msg.Message.Chat.ID, "usage example: 02.01.2006_15:04 once Worst idea for date parsing"))
		return
	}

	rm := getNewReminder(date.UTC(), settings.User, msg.Message.Text, parseInterval(i), msg.Message.Chat.ID)

	err = mongoCtl.addReminder(rm)

	if err != nil {
		bot.Send(tgbotapi.NewMessage(msg.Message.Chat.ID, "failed to add reminder"))
		return
	}

	bot.Send(tgbotapi.NewMessage(msg.Message.Chat.ID, "registered reminder"))
}
