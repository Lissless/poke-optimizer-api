package main

import (
	"log"
	"net/http"

	pc "github.com/Lissless/poke-optimizer-api/pc_ops"
	"github.com/Lissless/poke-optimizer-api/pkmn"
	"github.com/rs/cors"
)

func main() {
	r := http.NewServeMux()

	r.Handle("/pokemon/", &pkmn.PokemonHandler{})
	r.Handle("/pc/", &pc.PCHandler{
		ActiveTeam: pc.PokemonTeam{
			PokemonTeam: make([]pkmn.ActivePokemon, 0),
		},
		PkmnBox: pc.PokemonStorageBox{
			PCBox: make([]pkmn.ActivePokemon, 0),
		},
	})

	// Solves Cross Origin Access Issue
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	handler := c.Handler(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":8080", // The port to listen into
	}

	log.Println("Pokemon Optimizer has started listening")

	log.Fatal(srv.ListenAndServe())
}
