package main

import (
	"log"
	"net/http"

	"github.com/Lissless/poke-optimizer-api/pkmn"
	"github.com/rs/cors"
)

func main() {
	r := http.NewServeMux()

	r.Handle("/pokemon/", &pkmn.PokemonHandler{})

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
