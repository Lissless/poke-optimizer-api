package pc

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Lissless/poke-optimizer-api/pkmn"
	"github.com/Lissless/poke-optimizer-api/pkmn_errors"
)

type PCHandler struct {
	ActiveTeam PokemonTeam
	PkmnBox    PokemonStorageBox
}

func (pc *PCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		pc.GetPokemonTeam(w, r)
		return
	case r.Method == http.MethodPut:
		pc.AppendToPokemonTeam(w, r)
		return
	default:
		log.Printf("Invalid route was attempted route: %s", r.URL)
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Invalid request")
	}
}

func (pc *PCHandler) GetPokemonTeam(w http.ResponseWriter, r *http.Request) {
	write_resp, err := json.Marshal(pc.ActiveTeam)
	if err != nil {
		log.Printf("Marshalling the request to get a pokemon team: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon team request")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}

func (pc *PCHandler) AppendToPokemonTeam(w http.ResponseWriter, r *http.Request) {
	if len(pc.ActiveTeam.PokemonTeam) == 6 {
		// cannot have more than 6 pokemon in a team
		log.Printf("Attempted to add more pokrmon to a full team")
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed: current team is too full")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Reading request body to GET the url: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading request body")
		return
	}
	defer r.Body.Close()

	newPKMN := pkmn.ActivePokemon{}
	err = json.Unmarshal(body, &newPKMN)
	// change this later
	if err != nil {
		log.Printf("Reading request body to GET the url: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading request body")
		return
	}

	pc.ActiveTeam.PokemonTeam = append(pc.ActiveTeam.PokemonTeam, newPKMN)

	// Todo: return the new pc team immediately

}
