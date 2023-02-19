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

type App struct {
	Router    *mux.Router
	DataStore *DataStore
	Upgrader  websocket.Upgrader
	Database  *mongo.Database
}

func (app *App) Initialize() {
	db, err := connectDb()
	if err != nil {
		log.Println("Failed to connect mongodb")
	}

	app.Router = mux.NewRouter()
	app.DataStore = &DataStore{data: make(map[string]interface{})}
	app.Upgrader = websocket.Upgrader{}
	app.Database = db

	app.InitializeRoutes()
}

func (app *App) Run(addr string) {
	log.Printf("Webserver is started at port: %s\n", addr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", addr), app.Router))
}

func connectDb() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://mongodb:27017")
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
		log.Println("err ping", err)
		return nil, err
	}

	log.Println("connected to mongodb")

	return client.Database("key_value_store"), nil
}

func main() {
	app := App{}

	app.Initialize()
	app.Run(webPort)
}
