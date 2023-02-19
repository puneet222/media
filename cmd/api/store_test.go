package main

import "testing"

func TestDataStore_AddAndGet(t *testing.T) {
	store := DataStore{data: make(map[string]interface{})}

	// Add a key-value pair to the data store
	store.Add("key1", "value1")

	// Get the value for an existing key
	value, ok := store.Get("key1")
	if !ok {
		t.Errorf("Expected to find key 'key1', but didn't")
	}
	if value != "value1" {
		t.Errorf("Expected value for key 'key1' to be 'value1', but got '%v'", value)
	}

	// Get the value for a non-existent key
	value, ok = store.Get("key2")
	if ok {
		t.Errorf("Expected not to find key 'key2', but did")
	}
	if value != nil {
		t.Errorf("Expected value for non-existent key to be nil, but got '%v'", value)
	}
}
