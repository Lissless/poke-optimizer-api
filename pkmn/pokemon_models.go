package pkmn

// there is a "past types" field which is game specific
// add functinoality when looking at a particular game
type Pokemon struct {
	Name      string           `json:"name"`
	ID        float64          `json:"id"`
	Types     []PokemonType    `json:"types"`
	Abilities []PokemonAbility `json:"abilities"`
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
