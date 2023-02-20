package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

const webPort = "8088"
const dbName = "keyValueDB"
const collection = "keyValue"

type App struct {
	Router     *mux.Router
	DataStore  *DataStore
	Upgrader   websocket.Upgrader
	DB         *mongo.Database
	DbName     string
	Collection string
	ctx        context.Context
}

func (app *App) Initialize() {
	db, err := app.connectDb(dbName)
	if err != nil {
		log.Println("Failed to connect mongodb")
	}
	log.Println("connected to mongodb")

	app.Router = mux.NewRouter()
	app.DataStore = &DataStore{data: make(map[string]interface{})}
	app.Upgrader = websocket.Upgrader{}
	app.DB = db
	app.DbName = dbName
	app.Collection = collection
	app.ctx = context.Background()

	// update datastore from database
	app.UpdateDataStore()

	app.InitializeRoutes()
}

func (app *App) Run(addr string) {
	log.Printf("Webserver is started at port: %s\n", addr)
	err := http.ListenAndServe(fmt.Sprintf(":%s", addr), app.Router)
	if err != nil {
		log.Fatal("Error while starting webserver", err)
	}
}

func (app *App) UpdateDataStore() {
	keyValues := app.GetAllKeyValues()
	for _, kv := range keyValues {
		app.DataStore.Add(kv.Key, kv.Value)
	}
}

func main() {
	app := App{}

	app.Initialize()
	app.Run(webPort)
}
