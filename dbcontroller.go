package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoopts "go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURI = "mongodb://localhost:27017"

type dbController struct {
	client     *mongo.Client
	clientOpts *mongoopts.ClientOptions
	colls      map[string]*mongo.Collection
}

func (db *dbController) Init() {
	db.clientOpts = mongoopts.Client().ApplyURI(mongoURI)

	var err error

	db.client, err = mongo.Connect(context.TODO(), db.clientOpts)

	if err != nil {
		log.Panic("failed to connect to mongo", err)
		return
	}

	err = db.client.Ping(context.TODO(), nil)

	if err != nil {
		log.Panic("failed to ping", err)
		return
	}

	db.colls = make(map[string]*mongo.Collection)

	db.colls["reminders"] = db.client.Database("test").Collection("reminders")
	db.colls["users"] = db.client.Database("test").Collection("users")
}

func (db *dbController) getCollection(name string) *mongo.Collection {
	return db.colls[name]
}

func (db *dbController) Cleanup() {
	db.client.Disconnect(context.TODO())
}

func (db *dbController) findUserSettings(name string) (UserSettings, error) {
	filter := bson.D{{Key: "username", Value: name}}

	var u UserSettings

	e := db.colls["users"].FindOne(context.TODO(), filter).Decode(&u)

	return u, e
}

func (db *dbController) addUserSettings(u UserSettings) error {
	_, err := db.colls["users"].InsertOne(context.TODO(), u)

	return err
}

func (db *dbController) removeUserSettings(name string) error {
	filter := bson.D{{Key: "username", Value: name}}

	_, err := db.colls["users"].DeleteOne(context.TODO(), filter)

	return err
}

func (db *dbController) addReminder(r Reminder) error {
	_, err := db.colls["reminders"].InsertOne(context.TODO(), r)

	return err
}

func (db *dbController) removeReminders(key, value string) (int64, error) {
	filter := bson.D{{Key: key, Value: value}}

	res, err := db.colls["reminders"].DeleteMany(context.TODO(), filter)

	return res.DeletedCount, err
}

func (db *dbController) modifyReminder(key, value string, newValue Reminder) {
	db.colls["reminders"].FindOneAndReplace(context.TODO(), bson.D{{Key: key, Value: value}}, newValue)
}

func (db *dbController) findReminders(key, value string) []Reminder {
	filter := bson.D{{Key: key, Value: value}}

	cur, err := db.colls["reminders"].Find(context.TODO(), filter)

	res := make([]Reminder, 0)

	if err != nil {
		return res
	}

	for cur.Next(context.TODO()) {
		rm := Reminder{}

		err = cur.Decode(&rm)

		if err != nil {
			continue
		}

		res = append(res, rm)
	}

	return res
}
