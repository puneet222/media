package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func (app *App) InsertToDB(data KeyValue) {
	_, err := app.DB.Collection(app.Collection).InsertOne(app.ctx, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (app *App) GetAllKeyValues() []KeyValue {
	cur, err := app.DB.Collection(app.Collection).Find(app.ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		err := cur.Close(app.ctx)
		if err != nil {
			log.Println("Error while closing mongo collection cursor", err)
		}
	}()

	result := make([]KeyValue, 0)
	for cur.Next(app.ctx) {
		var row KeyValue
		err := cur.Decode(&row)
		if err != nil {
			log.Fatal("error while decoding mongo row", err.Error())
		}

		result = append(result, row)
	}
	return result
}
