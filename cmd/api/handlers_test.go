package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApp_GetData(t *testing.T) {
	store := DataStore{data: make(map[string]interface{})}
	store.Add("key1", "value1")
	app := App{DataStore: &store}

	// Test retrieving an existing key
	req, err := http.NewRequest("GET", "/data?key=key1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetData)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}
	expected := KeyValue{Key: "key1", Value: "value1"}
	var actual KeyValue
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v, expected %v", actual, expected)
	}

	// Test retrieving a non-existent key
	req, err = http.NewRequest("GET", "/data?key=baz", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v, expected %v", status, http.StatusNotFound)
	}
	expectedError := "Key not found"

	if actualError := strings.TrimSuffix(rr.Body.String(), "\n"); actualError != expectedError {
		t.Errorf("Handler returned unexpected error message: got '%v', expected '%v'", actualError, expectedError)
	}

	// Test missing key parameter
	req, err = http.NewRequest("GET", "/data", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v, expected %v", status, http.StatusBadRequest)
	}
	expectedError = "Missing key parameter"
	if actualError := strings.TrimSuffix(rr.Body.String(), "\n"); actualError != expectedError {
		t.Errorf("Handler returned unexpected error message: got %v, expected %v", actualError, expectedError)
	}
}

func TestApp_createWebsocket(t *testing.T) {
	store := DataStore{data: make(map[string]interface{})}
	app := App{DataStore: &store}

	// Create a new WebSocket connection
	server := httptest.NewServer(http.HandlerFunc(app.createWebsocket))
	defer server.Close()
	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/"
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error while closing the websocket connection")
		}
	}()

	// Test sending a key-value pair over the WebSocket
	keyValue := KeyValue{Key: "key1", Value: "value1"}
	m, err := json.Marshal(keyValue)
	if err != nil {
		t.Fatal(err)
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, m); err != nil {
		t.Fatal(err)
	}

	var resp WSResponse
	err = conn.ReadJSON(&resp)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the key-value pair was stored in the data store
	value, ok := store.Get("key1")
	if !ok {
		t.Errorf("Expected to find key 'key1', but didn't")
	}
	if value != "value1" {
		t.Errorf("Expected to find key 'value1', but didn't")
	}

	// Check that we got the expected response
	//if resp.success == false {
	//	t.Errorf("Expected to find success value to true but got false")
	//}
	//if resp.message != "Key value pair added successfully" {
	//	t.Errorf("Got unexpected response message [%s]", resp.message)
	//}
}
