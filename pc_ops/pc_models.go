package pc

import "github.com/Lissless/poke-optimizer-api/pkmn"

// Every user will have their own "active team" and their own
// storage box, find a way to store this per person, mongo?

type PokemonTeam struct {
	PokemonTeam []pkmn.ActivePokemon `json:"pokemon_team"`
}

type PokemonStorageBox struct {
	PCBox []pkmn.ActivePokemon `json:"pc_box_pokemon"`
}
