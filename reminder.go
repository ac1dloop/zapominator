package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//TimeInterval will be parsed according to interval
type TimeInterval int

const (
	once TimeInterval = 0
	day
	month
	week
	year
	custom
)

//Reminder contains data used to restore message
//then send it to correct chat
type Reminder struct {
	ID       int64         `json:"id"`
	Date     time.Time     `json:"date"`
	User     tgbotapi.User `json:"user"`
	Text     string        `json:"text"`
	Interval TimeInterval  `json:"interval"`
	ChatID   int64         `json:"chat_id"`
}

//ReminderID Global ID for reminders. incremented each time new is created
var ReminderID int64 = 0

func getNewReminder(Date time.Time, User tgbotapi.User, Text string, Interval TimeInterval, ChatID int64) Reminder {
	ReminderID++

	return Reminder{ReminderID - 1, Date, User, Text, Interval, ChatID}
}
