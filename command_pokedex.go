package main

import (
	"fmt"
)

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.caughtPokemon) == 0 {
		fmt.Println("you have not caught any pokemon")
		return nil
	}
	for name, _ := range cfg.caughtPokemon {
		fmt.Printf("- %s\n", name)
	}
	return nil
}
