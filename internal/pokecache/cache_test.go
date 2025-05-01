package pokecache

import (
	"testing"
	"time"
)

func TestExpiredEntryIsDeleted(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	cache.Add("expired-key", []byte("old"))

	// simulate manual expiration by manipulating the internal timestamp
	entry := cacheEntry{
		val:       []byte("old"),
		createdAt: time.Now().Add(-10 * time.Minute), // expired
	}
	cache.entries["expired-key"] = entry

	val, ok := cache.Get("expired-key")
	if ok || val != nil {
		t.Errorf("expected expired key to be deleted and not returned")
	}
}

func TestNonExistentKeyReturnsFalse(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	val, ok := cache.Get("missing-key")
	if ok || val != nil {
		t.Errorf("expected missing key to return false and nil value")
	}
}

func TestGetRefreshesDoesNotRefreshTime(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	cache.Add("key", []byte("data"))
	firstTime := cache.entries["key"].createdAt

	time.Sleep(10 * time.Millisecond)
	_, _ = cache.Get("key")

	if cache.entries["key"].createdAt != firstTime {
		t.Errorf("expected Get() not to refresh createdAt timestamp")
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	go cache.reapLoop(baseTime)

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
	}
}
