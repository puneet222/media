package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type KeyValue struct {
	Key   string `json:"key" bson:"key"`
	Value any    `json:"value" bson:"value"`
}

type WSResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (app *App) GetData(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	value, ok := app.DataStore.Get(key)
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(KeyValue{Key: key, Value: fmt.Sprint(value)})
	if err != nil {
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Error while writing response", http.StatusInternalServerError)
		return
	}
}

func (app *App) createWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := app.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error upgrading to WebSocket protocol", http.StatusInternalServerError)
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal("Error while closing websocket connection")
		}
	}()

	for {
		var keyValue KeyValue
		err := conn.ReadJSON(&keyValue)
		if err != nil {
			break
		}
		// Add data to in memory datastore
		app.DataStore.Add(keyValue.Key, keyValue.Value)

		// Add data to persistent storage in a new go routine
		go app.InsertToDB(keyValue)

		resp := WSResponse{
			Success: true,
			Message: "Key value pair added successfully",
		}
		m, err := json.Marshal(resp)
		if err != nil {
			log.Println("Error while marshaling json response", err)
		}
		err = conn.WriteMessage(websocket.BinaryMessage, m)
		if err != nil {
			log.Println(err)
		}
	}
}
