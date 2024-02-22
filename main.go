package main

import (
	"fmt"
	"github.com/Lissless/poke-optimizer-api/pkmn"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", pkmn.GetPokemon).Methods("GET")

	// Solves Cross Origin Access Issue
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	handler := c.Handler(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":8080", // The port to listen into
	}

	log.Println("Ok we about to listen")

	log.Fatal(srv.ListenAndServe())
}
