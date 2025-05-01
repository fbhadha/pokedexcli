package main

import (
	"fmt"
	"os"
)

func commandExit(cfg *config, params []string) error {
	// if len(params) > 0 {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	// }
	return nil
}
