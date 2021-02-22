package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoCtl = dbController{}
var bot = tgBot{}

func main() {

	vers := flag.Bool("v", false, "print version")

	flag.Parse()

	if *vers {
		fmt.Println(version + " " + versionComment)

		return
	}

	mongoCtl.Init()
	bot.Init()

	bot.addHandler("/remind", setReminderHandler)
	bot.addHandler("/help", helpHandler)
	bot.addHandler("/test", testHandler)
	bot.addHandler("/start", startHandler)
	bot.addHandler("/settings", settingsHandler)
	bot.addHandler("/ping", pingHandler)

	go func() {
		for {
			checkReminders(mongoCtl.getCollection("reminders"), bot.ptr)
		}
	}()

	log.Println("Bot started")
	bot.Start()
}

//infinite check for reminder
func checkReminders(col *mongo.Collection, bot *tgbotapi.BotAPI) {

	for {

		cur, e := col.Find(context.TODO(), bson.D{})

		if e != nil {
			log.Fatal(e)
			continue
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			rm := Reminder{}

			e = cur.Decode(&rm)

			if e != nil {
				log.Fatal(e)
				continue
			}

			//delete object if its outdated
			if time.Until(rm.Date) < 0 {

				col.DeleteOne(context.TODO(), cur.Current)

				continue
			}

			if time.Until(rm.Date) <= time.Minute {

				m := tgbotapi.NewMessage(rm.ChatID, rm.Text)

				_, e = bot.Send(m)

				if e != nil {
					log.Panic("failed to send message: ", e)
				}

				if rm.Interval == once {
					col.DeleteOne(context.TODO(), cur.Current)
				}
			}
		}

		time.Sleep(time.Second)

	}

}
