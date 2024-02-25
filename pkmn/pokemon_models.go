package pkmn

import (
	"crypto/rand"
	"encoding/hex"
)

// there is a "past types" field which is game specific
// add functinoality when looking at a particular game
type Pokemon struct {
	Name      string           `json:"name"`
	PokedexID float64          `json:"id"`
	Types     []PokemonType    `json:"types"`
	Abilities []PokemonAbility `json:"abilities"`
}

// this represents a pokemon that is "owned" by the user,
type ActivePokemon struct {
	UID       string         `json:"uid"` // so even if there are multiple of the same pokemon specific ones can be identified
	Name      string         `json:"name"`
	PokedexID float64        `json:"id"`
	Types     []PokemonType  `json:"types"`
	Abilities PokemonAbility `json:"ability"`
}

type PokemonType struct {
	Name string  `json:"type_name"`
	Slot float64 `json:"slot"`
}

type PokemonAbility struct {
	Ability   string  `json:"ability_name"`
	Is_hidden bool    `json:"is_hidden"`
	Slot      float64 `json:"slot"`
}

func (ap *ActivePokemon) GeneratePkmnUID() error {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	ap.UID = hex.EncodeToString(bytes)
	return nil
}
