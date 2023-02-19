package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func (app *App) Insert(value KeyValue) {
	var ctx = context.Background()
	_, err := app.Database.Collection("data").InsertOne(ctx, bson.D{{value.Key, value.Value}})
	if err != nil {
		log.Fatal(err.Error())
	}
}
