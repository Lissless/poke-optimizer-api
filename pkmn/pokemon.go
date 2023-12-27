package pkmn

import (
	"log"
	"net/http"
)

func GetPokemon(w http.ResponseWriter, r *http.Request) {
	log.Println("We did it!")
}
