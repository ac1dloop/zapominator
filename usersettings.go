package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//UserSettings contains data with
//individual user settings
type UserSettings struct {
	User     tgbotapi.User
	Location string `json:"location"`
	Layout   string `json:"layout"`
}

//returns default settings if specific not found
func getDefaultUserSettings() UserSettings {
	return UserSettings{tgbotapi.User{}, "Europe/Moscow", "02.01.2006_15:04"}
}
