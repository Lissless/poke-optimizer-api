package pkmn

func MakePokemon(dataMap map[string]interface{}, name string) Pokemon {
	pkmn := Pokemon{
		Name:      name,
		PokedexID: dataMap["id"].(float64),
		Types:     makeTypesArray(dataMap["types"].([]interface{})),
		Abilities: makeAbilitiesArray(dataMap["abilities"].([]interface{})),
	}

	return pkmn
}

func makeTypesArray(raw_types []interface{}) []PokemonType {
	pokemon_types := []PokemonType{}
	for _, data := range raw_types {
		type_data := data.(map[string]interface{})
		pokemon_types = append(pokemon_types, PokemonType{
			Slot: type_data["slot"].(float64),
			Name: type_data["type"].(map[string]interface{})["name"].(string),
		})
	}
	return pokemon_types
}

func makeAbilitiesArray(raw_abilities []interface{}) []PokemonAbility {
	pokemon_abilities := []PokemonAbility{}
	for _, data := range raw_abilities {
		ability_data := data.(map[string]interface{})
		pokemon_abilities = append(pokemon_abilities, PokemonAbility{
			Ability:   ability_data["ability"].(map[string]interface{})["name"].(string),
			Is_hidden: ability_data["is_hidden"].(bool),
			Slot:      ability_data["slot"].(float64),
		})
	}
	return pokemon_abilities
}
