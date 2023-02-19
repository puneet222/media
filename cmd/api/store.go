package main

type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type DataStore struct {
	data map[string]interface{}
}

func (store *DataStore) Add(key string, value interface{}) {
	store.data[key] = value
}

func (store *DataStore) Get(key string) (interface{}, bool) {
	value, ok := store.data[key]
	return value, ok
}
