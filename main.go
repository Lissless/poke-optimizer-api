package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"poke-optimizer-api/pkmn"
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
		Addr:    ":80", // The port to listen into
	}

	log.Println("Ok we about to listen")

	log.Fatal(srv.ListenAndServe())
}

// func getPokemon(w http.ResponseWriter, r *http.Request) {
// 	log.Println("We did it!")
// }
