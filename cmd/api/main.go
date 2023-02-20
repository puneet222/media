package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

const webPort = "8088"
const mongoDBUri = "mongodb://mongodb:27017"
const dbName = "keyValueDB"
const collection = "keyValue"

type App struct {
	Router     *mux.Router
	DataStore  *DataStore
	Upgrader   websocket.Upgrader
	Database   *mongo.Database
	DbName     string
	Collection string
}

func (app *App) Initialize() {
	db, err := connectDb(dbName)
	if err != nil {
		log.Println("Failed to connect mongodb")
	}
	log.Println("connected to mongodb")

	app.Router = mux.NewRouter()
	app.DataStore = &DataStore{data: make(map[string]interface{})}
	app.Upgrader = websocket.Upgrader{}
	app.Database = db
	app.DbName = dbName
	app.Collection = collection

	app.InitializeRoutes()
}

func (app *App) Run(addr string) {
	log.Printf("Webserver is started at port: %s\n", addr)
	err := http.ListenAndServe(fmt.Sprintf(":%s", addr), app.Router)
	if err != nil {
		log.Fatal("Error while starting webserver", err)
	}
}

func connectDb(databaseName string) (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(mongoDBUri)
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
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

func main() {
	app := App{}

	app.Initialize()
	app.Run(webPort)
}
