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
	RELEASE_PATH     = "/pc/release"
	POKEMON_PATH     = "/pc/pokemon"
)

func (pc *PCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && (r.URL.Path == ACTIVE_TEAM_PATH || r.URL.Path == BOX_PATH):
		pc.GetPokemonStorageLocation(w, r)
		return
	case r.Method == http.MethodGet && r.URL.Path == POKEMON_PATH:
		pc.GetSpecificPokemon(w, r)
		return
	case r.Method == http.MethodPut && r.URL.Path == ACTIVE_TEAM_PATH:
		pc.AddPokemon(w, r, true)
		return
	case r.Method == http.MethodPut && r.URL.Path == BOX_PATH:
		pc.AddPokemon(w, r, false)
		return
	case r.Method == http.MethodPut && r.URL.Path == TRANSFER_PATH:
		pc.TransferPokemon(w, r)
		return
	case r.Method == http.MethodDelete && r.URL.Path == RELEASE_PATH:
		pc.ReleasePokemon(w, r)
		return
	default:
		log.Printf("Invalid route was attempted route: %s", r.URL)
		pkmn_errors.ErrorHandler(w, r, http.StatusBadRequest, "Invalid request")
	}
}

func (pc *PCHandler) GetPokemonStorageLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Get Pokemon location info for storage info: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading get information request body")
		return
	}
	defer r.Body.Close()

	pkmnLocation := PokemonLocation{}
	err = json.Unmarshal(body, &pkmnLocation)
	if err != nil {
		log.Printf("Unmarshalling the Pokemon location for storage info failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading get information request body")
		return
	}

	var write_resp []byte
	switch pkmnLocation.Location {
	case "team":
		write_resp, err = json.Marshal(pc.ActiveTeam)
	case "box":
		write_resp, err = json.Marshal(pc.PkmnBox)
	}

	if err != nil {
		log.Printf("Marshalling the request to get the pokemon %s failed, error: %s", pkmnLocation.Location, err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon box request")
		return
	}

	log.Printf("Retrieved information on the pokemon %s", pkmnLocation.Location)
	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}

func (pc *PCHandler) AddPokemon(w http.ResponseWriter, r *http.Request, locTeam bool) {
	if locTeam {
		// cannot have more than 6 pokemon in a team
		if !EnforceTeamLimit(pc.ActiveTeam) {
			pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed: current team is too full")
			return
		}
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
	var locationString string
	if locTeam {
		pc.ActiveTeam.PokemonTeam[newPKMN.UID] = newPKMN
		write_resp, err = json.Marshal(pc.ActiveTeam.PokemonTeam)
		locationString = "team"
	} else {
		pc.PkmnBox.PCBox[newPKMN.UID] = newPKMN
		write_resp, err = json.Marshal(pc.PkmnBox.PCBox)
		locationString = "box"
	}

	if err != nil {
		log.Printf("Marshalling the request to GET the pokemon team or storage, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon team or storage request")
		return
	}

	log.Println("Successfully added a new pokemon to the " + locationString)
	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}

func (pc *PCHandler) TransferPokemon(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Transfer Pokemon Reading request body to GET the url: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading transfer request body")
		return
	}
	defer r.Body.Close()

	pkmnLocation := PokemonLocation{}
	err = json.Unmarshal(body, &pkmnLocation)
	if err != nil {
		log.Printf("Transfer request body to GET the pokemon for adding to the team or storage: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading transfer request body")
		return
	}

	var write_resp []byte
	valid := false
	switch pkmnLocation.Location {
	case "team":
		pc.ActiveTeam, pc.PkmnBox, valid = MovePokemonFromTeamToBox(pc.ActiveTeam, pc.PkmnBox, pkmnLocation.PokemonUID)
		if !valid {
			pkmn_errors.ErrorHandler(w, r, http.StatusBadRequest, "Failed to get pokemon from the active team")
			return
		}
		write_resp, err = json.Marshal(pc.ActiveTeam)
	case "box":
		pc.ActiveTeam, pc.PkmnBox, valid = MovePokemonFromBoxToTeam(pc.ActiveTeam, pc.PkmnBox, pkmnLocation.PokemonUID)
		if !valid {
			pkmn_errors.ErrorHandler(w, r, http.StatusBadRequest, "Failed to get pokemon from the pc box")
			return
		}
		write_resp, err = json.Marshal(pc.PkmnBox)
	default:
		log.Printf("Invalid location for transfer was attempted, route: %s", r.URL)
		pkmn_errors.ErrorHandler(w, r, http.StatusBadRequest, "Invalid transfer location")
		return
	}

	if err != nil {
		log.Printf("Marshalling the request to transfer pokemon, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon transfer request")
		return
	}

	log.Println("Sucessfully transfered pokemon from: " + pkmnLocation.Location)
	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}

func (pc *PCHandler) ReleasePokemon(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Transfer Pokemon Reading request body to GET the url: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading transfer request body")
		return
	}
	defer r.Body.Close()

	pkmnLocation := PokemonLocation{}
	err = json.Unmarshal(body, &pkmnLocation)
	if err != nil {
		log.Printf("Transfer request body to GET the pokemon for adding to the team or storage: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading transfer request body")
		return
	}

	var write_resp []byte
	valid := false
	switch pkmnLocation.Location {
	case "team":
		_, valid = pc.ActiveTeam.RemovePokemon(pkmnLocation.PokemonUID)
		write_resp, err = json.Marshal(pc.ActiveTeam)
	case "box":
		_, valid = pc.PkmnBox.RemovePokemon(pkmnLocation.PokemonUID)
		write_resp, err = json.Marshal(pc.PkmnBox)
	}

	if !valid {
		log.Printf("Failed to retrieve and remove a pokemon from the %s, pokemon uid %s", pkmnLocation.Location, pkmnLocation.PokemonUID)
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to release the desired pokemon")
		return
	}

	if err != nil {
		log.Printf("Marshalling the request to delete pokemon, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon release request")
		return
	}

	log.Printf("Successfully removed a pokemon from the %s, uid: %s", pkmnLocation.Location, pkmnLocation.PokemonUID)
	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}

func (pc *PCHandler) GetSpecificPokemon(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Get Pokemon location info for storage info: failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading get information request body")
		return
	}
	defer r.Body.Close()

	pkmnLocation := PokemonLocation{}
	err = json.Unmarshal(body, &pkmnLocation)
	if err != nil {
		log.Printf("Unmarshalling the Pokemon location for storage info failed, error: %s", err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed reading get information request body")
		return
	}

	var write_resp []byte
	var pkmn pkmn.ActivePokemon
	valid := false
	switch pkmnLocation.Location {
	case "team":
		pkmn, valid = pc.ActiveTeam.GetPokemon(pkmnLocation.PokemonUID)
		write_resp, err = json.Marshal(pkmn)
	case "box":
		pkmn, valid = pc.PkmnBox.GetPokemon(pkmnLocation.PokemonUID)
		write_resp, err = json.Marshal(pkmn)
	}

	if !valid {
		log.Printf("Failed to retrieve a pokemon from the %s, pokemon uid %s", pkmnLocation.Location, pkmnLocation.PokemonUID)
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to retrieve the desired pokemon")
		return
	}

	if err != nil {
		log.Printf("Marshalling the request to retrieve pokemon from the %s, error: %s", pkmnLocation.Location, err.Error())
		pkmn_errors.ErrorHandler(w, r, http.StatusInternalServerError, "Failed packaging pokemon pc retrieval request")
		return
	}

	log.Printf("Successfully retrieved a pokemon from the %s, uid: %s", pkmnLocation.Location, pkmnLocation.PokemonUID)
	w.WriteHeader(http.StatusOK)
	w.Write(write_resp)
}
