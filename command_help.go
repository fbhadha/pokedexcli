package main

import "fmt"

func commandHelp(cfg *config, param []string) error {
	// if len(param) > 0 {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	// }
	return nil
}
