package pc

import (
	"log"

	"github.com/Lissless/poke-optimizer-api/pkmn"
)

func MovePokemonFromTeamToBox(pkmnTeam PokemonTeam, pkmnBox PokemonStorageBox, pkmnUID string) (PokemonTeam, PokemonStorageBox, bool) {
	pkmn, valid := pkmnTeam.RemovePokemon(pkmnUID)
	if !valid {
		return pkmnTeam, pkmnBox, false
	}
	// add pokemon to the box
	pkmnBox.PCBox[pkmnUID] = pkmn
	return pkmnTeam, pkmnBox, true
}

func MovePokemonFromBoxToTeam(pkmnTeam PokemonTeam, pkmnBox PokemonStorageBox, pkmnUID string) (PokemonTeam, PokemonStorageBox, bool) {
	pkmn, valid := pkmnBox.RemovePokemon(pkmnUID)
	if !valid {
		return pkmnTeam, pkmnBox, false
	}
	if !EnforceTeamLimit(pkmnTeam) {
		return pkmnTeam, pkmnBox, false
	}
	// add pokemon to the team
	pkmnTeam.PokemonTeam[pkmnUID] = pkmn
	return pkmnTeam, pkmnBox, true
}

func EnforceTeamLimit(pkmnTeam PokemonTeam) bool {
	if len(pkmnTeam.PokemonTeam) == 6 {
		// cannot have more than 6 pokemon in a team
		log.Printf("Attempted to add more pokemon to a full team")
		return false
	}
	return true
}

func (pt *PokemonTeam) RemovePokemon(PkmnUID string) (pkmn.ActivePokemon, bool) {
	target := pt.PokemonTeam[PkmnUID]
	if target.UID == "" {
		return target, false
	}
	delete(pt.PokemonTeam, PkmnUID)
	return target, true
}

func (pb *PokemonStorageBox) RemovePokemon(PkmnUID string) (pkmn.ActivePokemon, bool) {
	target := pb.PCBox[PkmnUID]
	if target.UID == "" {
		return target, false
	}
	delete(pb.PCBox, PkmnUID)
	return target, true
}
