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

const (
	ACTIVE_TEAM_PATH = "/pc/team"
	BOX_PATH         = "/pc/box"
	TRANSFER_PATH    = "/pc/transfer"
)

func (pc *PCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == ACTIVE_TEAM_PATH:
		pc.GetPokemonTeam(w, r)
		return
	case r.Method == http.MethodGet && r.URL.Path == BOX_PATH:
		pc.GetPokemonBox(w, r)
		return
	case r.Method == http.MethodPut && r.URL.Path == ACTIVE_TEAM_PATH:
		pc.AddPokemon(w, r, true)
		return
	case r.Method == http.MethodPut && r.URL.Path == BOX_PATH:
		pc.AddPokemon(w, r, false)
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

func (pc *PCHandler) GetPokemonBox(w http.ResponseWriter, r *http.Request) {
	write_resp, err := json.Marshal(pc.PkmnBox)
	if err != nil {
		log.Printf("Marshalling the request to get the pokemon box failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon box request")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}

func (pc *PCHandler) AddPokemon(w http.ResponseWriter, r *http.Request, locTeam bool) {
	if len(pc.ActiveTeam.PokemonTeam) == 6 && locTeam {
		// cannot have more than 6 pokemon in a team
		log.Printf("Attempted to add more pokemon to a full team")
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed: current team is too full")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Append Pokemon Reading request body to GET the url: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading request body")
		return
	}
	defer r.Body.Close()

	newPKMN := pkmn.ActivePokemon{}
	err = json.Unmarshal(body, &newPKMN)
	if err != nil {
		log.Printf("Reading request body to GET the pokemon for adding to the team or storage: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading request body to add pokemon to team or storage")
		return
	}

	err = newPKMN.GeneratePkmnUID()
	if err != nil {
		log.Printf("Failed generating UID for a pokemon, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed assigning the pokemon an unique identifier")
		return
	}

	var write_resp []byte
	if locTeam {
	pc.ActiveTeam.PokemonTeam = append(pc.ActiveTeam.PokemonTeam, newPKMN)
		write_resp, err = json.Marshal(pc.ActiveTeam.PokemonTeam)

	} else {
		pc.PkmnBox.PCBox = append(pc.PkmnBox.PCBox, newPKMN)
		write_resp, err = json.Marshal(pc.PkmnBox.PCBox)
	}

	if err != nil {
		log.Printf("Marshalling the request to GET the pokemon team or storage, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon team or storage request")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}
