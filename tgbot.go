package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const token = "1452972533:AAGcLM0mZnlZsCc5XdspATeAwgPQ7uM3fZY"

type tgBotHandler func(msg *tgbotapi.Update, bot *tgbotapi.BotAPI)

type tgBot struct {
	ptr *tgbotapi.BotAPI

	currUser *UserSettings

	handlers map[string]tgBotHandler
}

func (bot *tgBot) addHandler(path string, handler tgBotHandler) {
	bot.handlers[path] = handler
}

func (bot *tgBot) messageLoop() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.ptr.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
		return
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		bot.processUpdate(&update)
	}
}

func (bot *tgBot) processUpdate(update *tgbotapi.Update) {
	cmd := extractFirst(&update.Message.Text)

	if handler, ok := bot.handlers[cmd]; ok {
		handler(update, bot.ptr)
	}
}

func (bot *tgBot) Init() {
	var err error

	bot.ptr, err = tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	bot.handlers = make(map[string]tgBotHandler)
}

//Start blocks thread of execution
func (bot *tgBot) Start() {
	bot.messageLoop()
}
