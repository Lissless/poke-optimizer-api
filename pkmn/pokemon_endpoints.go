package pkmn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const pokeApiURL = "https://pokeapi.co/api/v2/pokemon/"

// Todo: find a way to catcch errors and write custom responses for them
func GetPokemon(w http.ResponseWriter, r *http.Request) {
	urlArr := strings.Split(r.URL.Path, "/")
	pokemonName := strings.ToLower(urlArr[len(urlArr)-1])
	log.Printf("Getting data to create the Pokemon: %s", pokemonName)
	getPokemonURL := pokeApiURL + pokemonName
	req, err := http.NewRequest("GET", getPokemonURL, nil)
	if err != nil {
		log.Printf("Creating a request to GET the url: %s failed, error: %s", getPokemonURL, err.Error())
		return
	}

	client := &http.Client{}

	resp, _ := client.Do(req)
	// if err != nil {
	// 	return err
	// }
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	dataMap := make(map[string]interface{})

	err = json.Unmarshal(body, &dataMap)

	pkmn := MakePokemon(dataMap, pokemonName)

	fmt.Println("Pokemon id is {}", pkmn.ID)
	write_resp, err := json.Marshal(pkmn)
	w.Write(write_resp)
	fmt.Println("done")
}
