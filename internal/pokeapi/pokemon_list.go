package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func extractPokemonNames(data []byte) ([]string, error) {
	var parsed struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	names := []string{}
	for _, p := range parsed.PokemonEncounters {
		names = append(names, p.Pokemon.Name)
	}
	return names, nil
}

func (c *Client) GetLocationPokemon(area string) ([]string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", area)

	// Check cache first
	if data, ok := c.cache.Get(url); ok {
		return extractPokemonNames(data)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Cache response
	c.cache.Add(url, body)

	return extractPokemonNames(body)
}
