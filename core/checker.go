package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func getUpdate(hash string) (bool, error) {
	const baseURL string = "https://api.modrinth.com/v2"
	res, err := http.Get(fmt.Sprintf("%v/version_file/%v", baseURL, hash))
	if err != nil {
		fmt.Println("Error in getting the data from modrinth")
		log.Fatal(err)
		return false, err
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		fmt.Println(res.StatusCode)
		return false, nil
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API error: %s", res.Status)
	}

	var m ModrinthJSON
	e := json.Unmarshal(body, &m)

	if e != nil {
		log.Fatal(e)
	}

	return true, nil
}

// CheckFunInModrinth Checks if the mod is in the modrinth or not and fills mod[IsModrinth] value to boolean
func CheckFunInModrinth(mods *map[string]*Mod) *map[string]*Mod {

	for _, mod := range *mods {
		var wg sync.WaitGroup
		wg.Add(1)

		go func(mod *Mod) {
			defer wg.Done()

			// GetUpdate gives if mod is from modrinth or not
			fromModrinth, err := getUpdate(mod.Hash)
			if err != nil {
				log.Fatal(err)
			}
			mod.IsModrinth = fromModrinth
		}(mod)

		wg.Wait()
	}

	return mods
}
