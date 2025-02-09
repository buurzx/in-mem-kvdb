package inmemory

import (
	"testing"
)

func TestHashTableSet(t *testing.T) {
	ht := NewHashTable()
	ht.Set("key1", "value1")
	value, _ := ht.Get("key1")
	if value != "value1" {
		t.Errorf("expected value1, got %v", value)
	}
}

func TestHashTableGet(t *testing.T) {
	// ...existing code...
	ht := NewHashTable()
	ht.Set("key1", "value1")
	value, _ := ht.Get("key1")
	if value != "value1" {
		t.Errorf("expected value1, got %v", value)
	}
}

func TestHashTableDelete(t *testing.T) {
	ht := NewHashTable()
	ht.Set("key1", "value1")
	ht.Del("key1")
	value, exists := ht.Get("key1")
	if exists {
		t.Errorf("expected empty string, got %v", value)
	}
}
