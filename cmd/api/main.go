package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const webPort = "8088"

type App struct {
	Router    *mux.Router
	DataStore *DataStore
	Upgrader  websocket.Upgrader
}

func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	app.DataStore = &DataStore{data: make(map[string]interface{})}
	app.Upgrader = websocket.Upgrader{}

	app.InitializeRoutes()
}

func (app *App) Run(addr string) {
	log.Printf("Webserver is started at port: %s\n", addr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", addr), app.Router))
}

func main() {
	app := App{}
	app.Initialize()
	app.Run(webPort)
}
