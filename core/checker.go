package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"sync"
)

const baseURL string = "https://api.modrinth.com/v2"

// getFile helper function to extract the json response from the modrinth api
func getFile(mod *Mod) (bool, error) {
	res, err := http.Get(fmt.Sprintf("%v/version_file/%v", baseURL, mod.Hash))
	var modrinthJSON ModrinthJSON

	if err != nil {
		fmt.Println("Error in getting the data from modrinth")
		log.Fatal(err)
		return false, err
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API error: %s", res.Status)
	}

	e := json.Unmarshal(body, &modrinthJSON)

	if e != nil {
		log.Fatal(e)
	}
	mod.Files = modrinthJSON.Files
	mod.GameVersion = modrinthJSON.GameVersions
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
			fromModrinth, err := getFile(mod)
			if err != nil {
				log.Fatal(err)
			}
			mod.IsModrinth = fromModrinth

		}(mod)

		wg.Wait()

	}

	for _, mod := range *mods {
		update(mod)
	}
	return mods
}

// isUpdateNeeded helper function used to check if update is needed or not
func isUpdateNeeded(gameVersion []string, preferredVersion string) bool {
	return slices.Contains(gameVersion, preferredVersion)
}

func update(mod *Mod) {
	var downloadURL string
	var targetFileName string

	for _, file := range mod.Files {
		downloadURL = file.Url
		targetFileName = file.Filename
	}

	downloadMod(downloadURL, targetFileName)
}

func downloadMod(url string, name string) {

	if url == "" {
		fmt.Println("Empty URL detected")
		return
	}

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	fileName := "/Users/mayur/Downloads/" + name
	out, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Fatal(err)
	}

}
