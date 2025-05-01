package main

import (
	"fmt"
)

func commandExplore(cfg *config, param []string) error {
	if len(param) == 0 {
		return fmt.Errorf("please provide a location name")
	}

	locationName := param[0]
	pokemonResp, err := cfg.pokeapiClient.GetLocationPokemon(locationName)
	if err != nil {
		return err
	}

	if len(pokemonResp) == 0 {
		fmt.Println("No Pokémon found in this location.")
		return nil
	}

	fmt.Printf("Pokémon found in %s:\n", locationName)
	for _, name := range pokemonResp {
		fmt.Println(name)
	}
	return nil

}
