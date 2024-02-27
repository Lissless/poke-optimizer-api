package pkmn

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Lissless/poke-optimizer-api/pkmn_errors"
)

const pokeApiURL = "https://pokeapi.co/api/v2/pokemon/"

var (
	PokemonGet = regexp.MustCompile(`/pokemon/([[:alpha:]]+)$`)
)

type PokemonHandler struct{}

func (ph *PokemonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && PokemonGet.MatchString(r.URL.Path):
		ph.GetPokemon(w, r)
		return
	default:
		log.Printf("Invalid route was attempted route: %s", r.URL)
		pkmn_errors.ErrorHandler(w, r, http.StatusBadRequest, "Invalid request")
	}
}

func (ph *PokemonHandler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	urlArr := strings.Split(r.URL.Path, "/")
	pokemonName := strings.ToLower(urlArr[len(urlArr)-1])
	log.Printf("Getting data to create the Pokemon: %s", pokemonName)
	getPokemonURL := pokeApiURL + pokemonName

	req, err := http.NewRequest("GET", getPokemonURL, nil)
	if err != nil {
		log.Printf("Creating a request to GET the url: %s failed, error: %s", getPokemonURL, err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed creating request")
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Doing the request to GET the url: %s failed, error: %s", getPokemonURL, err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed executing request")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading request body to GET the url: %s failed, error: %s", getPokemonURL, err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading request body")
		return
	}
	dataMap := make(map[string]interface{})

	err = json.Unmarshal(body, &dataMap)
	if err != nil {
		log.Printf("Unmarshalling the request to GET the url: %s failed, error: %s", getPokemonURL, err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading request")
		return
	}

	pkmn := MakePokemon(dataMap, pokemonName)

	write_resp, err := json.Marshal(pkmn)
	if err != nil {
		log.Printf("Marshalling the request to GET the url: %s failed, error: %s", getPokemonURL, err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon request")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}
