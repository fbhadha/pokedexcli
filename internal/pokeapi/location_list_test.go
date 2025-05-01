// File: internal/pokeapi/location_list_test.go
package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fbhadha/pokedexcli/internal/pokecache"
)

func TestListLocations_ReturnsLocations(t *testing.T) {
	// Expected response struct
	expected := RespShallowLocations{
		Count:    1,
		Next:     nil,
		Previous: nil,
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{
				Name: "loc1",
				URL:  "http://example.com/loc1",
			},
		},
	}
	// Create a test server that returns valid JSON.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))
	defer ts.Close()

	// Create a client with test server's http client and a new cache.
	client := Client{
		httpClient: *ts.Client(),
		cache:      pokecache.NewCache(5 * time.Minute),
	}
	pageURL := ts.URL
	resp, err := client.ListLocations(&pageURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Count != expected.Count {
		t.Errorf("expected count %d, got %d", expected.Count, resp.Count)
	}
	if len(resp.Results) != len(expected.Results) {
		t.Errorf("expected %d results, got %d", len(expected.Results), len(resp.Results))
	}
	for i, res := range resp.Results {
		if res.Name != expected.Results[i].Name || res.URL != expected.Results[i].URL {
			t.Errorf("expected result %v, got %v", expected.Results[i], res)
		}
	}
}

func TestListLocations_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer ts.Close()

	client := Client{
		httpClient: *ts.Client(),
		cache:      pokecache.NewCache(5 * time.Minute),
	}
	pageURL := ts.URL
	_, err := client.ListLocations(&pageURL)
	if err == nil {
		t.Errorf("expected error due to invalid JSON but got nil")
	}
}

func TestListLocations_CacheHit(t *testing.T) {
	var requestCount int32

	// Create a test server with a counter.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		resp := RespShallowLocations{
			Count:    1,
			Next:     nil,
			Previous: nil,
			Results: []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				{
					Name: "loc-cache",
					URL:  "http://example.com/loc-cache",
				},
			},
		}
		b, err := json.Marshal(resp)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(b)
	}))
	defer ts.Close()

	client := Client{
		httpClient: *ts.Client(),
		cache:      pokecache.NewCache(5 * time.Minute),
	}
	pageURL := ts.URL

	// First call: should hit the server.
	_, err := client.ListLocations(&pageURL)
	if err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}

	// Second call: should be served from cache.
	_, err = client.ListLocations(&pageURL)
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}

	if atomic.LoadInt32(&requestCount) != 1 {
		t.Errorf("expected server to be hit once, got %d", requestCount)
	}
}
