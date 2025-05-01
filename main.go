package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type config struct {
	Next     string
	Previous string
}

var supportCommands map[string]cliCommand

func cleanInput(text string) []string {
	// Split the input string by spaces
	words := strings.Fields(text)
	// Create a slice to hold the cleaned words
	cleanedWords := make([]string, 0, len(words))
	// Iterate over each word
	for _, word := range words {
		// Trim leading and trailing spaces
		cleanedWord := strings.TrimSpace(word)
		loweredword := strings.ToLower(cleanedWord)
		// Append the cleaned word to the slice
		cleanedWords = append(cleanedWords, loweredword)
	}
	return cleanedWords
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex! \n Usage:")
	for name, cmd := range supportCommands {
		fmt.Printf("  %s: %s\n", name, cmd.description)
	}
	return nil
}

// Define struct for what you care about
type Pokedex struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
	GameIndices []struct {
		GameIndex  int `json:"game_index"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"game_indices"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}

type LocationList struct {
	Results []struct {
		Name string `json:"name"` // The name of the location
		URL  string `json:"url"`  // The URL to fetch more details about this location
	} `json:"results"`
	Next     string `json:"next"`     // The URL for the next page of locations (for pagination)
	Previous string `json:"previous"` // The URL for the previous page of locations (for pagination)
}

func commandMap(c *config) error {
	// Determine which URL to request:
	// If c.Next is set, it means the user has paged forward before, and we want to fetch the next page.
	// If c.Next is empty, we start from the first page of the locations endpoint.
	url := c.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location?limit=20" // Default to the first 20 locations
	}

	// Make an HTTP GET request to the chosen URL.
	// http.Get returns an http.Response and an error.
	// http.Get handles the underlying TCP connection, sending the request, and receiving the response.
	resp, err := http.Get(url)
	if err != nil {
		// If there was a network or protocol error, return it so the REPL can display it.
		return err
	}
	// Always close the response body when done reading to free up network resources.
	defer resp.Body.Close()

	// Read the entire response body into memory as a byte slice.
	// io.ReadAll reads all data from the response until EOF.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// If there was an error reading the response, return it.
		return err
	}

	// Unmarshal (decode) the JSON response into our LocationList struct.
	// json.Unmarshal parses the JSON bytes and populates the struct fields accordingly.
	var locations LocationList
	if err := json.Unmarshal(body, &locations); err != nil {
		// If the JSON is malformed or doesn't match our struct, return the error.
		return err
	}

	// Print each location in the current page of results.
	// We use a loop to enumerate and format each location name for user-friendly display.
	for i, loc := range locations.Results {
		fmt.Printf("%2d. %s\n", i+1, loc.Name) // Print the index (1-based) and the location name
	}

	// Update the config with the pagination URLs for "next" and "previous" pages.
	// This allows the user to navigate forward and backward through the list.
	c.Next = locations.Next
	c.Previous = locations.Previous

	// Return nil to indicate success to the REPL.
	return nil
}

func commandMapBack(c *config) error {
	// Use the previous page URL stored in the config to go back in the paginated results.
	url := c.Previous
	if url == "" {
		// If there is no previous page (we're at the beginning), inform the user.
		fmt.Println("You're at the beginning of the map.")
		return nil // Not an error; just can't go back further.
	}

	// Make an HTTP GET request to fetch the previous page of locations.
	// http.Get handles the network communication and returns a response or error.
	resp, err := http.Get(url)
	if err != nil {
		// If there was a network or protocol error, return it for the REPL to display.
		return err
	}
	// Always close the response body to free up resources.
	defer resp.Body.Close()

	// Read the entire response body into a byte slice.
	// io.ReadAll reads all the data until EOF.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// If there was an error reading the response, return it.
		return err
	}

	// Decode the JSON response into our LocationList struct.
	// json.Unmarshal parses the JSON and fills our struct fields.
	var locations LocationList
	if err := json.Unmarshal(body, &locations); err != nil {
		// If the JSON doesn't match our struct or is invalid, return the error.
		return err
	}

	// Print each location in the current (previous) page for the user to see.
	for i, loc := range locations.Results {
		fmt.Printf("%2d. %s\n", i+1, loc.Name) // Nicely format the index and location name
	}

	// Update the config's pagination URLs so the user can continue to navigate forward/backward.
	c.Next = locations.Next
	c.Previous = locations.Previous

	// Return nil to indicate success to the REPL.
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func initCommands() {
	supportCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back 20 locations",
			callback:    commandMapBack,
		},
	}
}

func main() {
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}
		cmdName := words[0]
		if cmd, ok := supportCommands[words[0]]; ok {
			if err := cmd.callback(cfg); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", cmdName)
		}
	}
}
