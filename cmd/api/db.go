package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const mongoDBUri = "mongodb://mongodb:27017"
const dbTimeout = time.Second * 10

func (app *App) connectDb(databaseName string) (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(mongoDBUri)
	var ctx, _ = context.WithTimeout(context.Background(), dbTimeout)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("error while pinging mongodb", err)
		return nil, err
	}

	return client.Database(databaseName), nil
}

func (app *App) InsertToDB(data KeyValue) {
	// using mutex to resolve concurrent
	app.dbMutex.Lock()
	defer app.dbMutex.Unlock()
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
