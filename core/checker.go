package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const baseURL = "https://api.modrinth.com/v2"

func GetModIds(mods map[string]*Mod) error {
	CheckIfInModrinth(mods)

	url := fmt.Sprintf("%s/version_files", baseURL)
	// url := "https://webhook.site/536486f9-f38c-4681-a555-71e168a233dc"

	hashes := make([]string, 0, len(mods))

	for _, mod := range mods {
		if mod.IsModrinth {
			hashes = append(hashes, mod.Hash)
		}
	}

	jsonData, err := json.Marshal(struct {
		Hashes    []string `json:"hashes"`
		Algorithm string   `json:"algorithm"`
	}{
		Hashes:    hashes,
		Algorithm: "sha1",
	})
	if err != nil {
		fmt.Println("Error in parsing body to JSON")
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if res.StatusCode == http.StatusNotFound {
		fmt.Println("This hash is not possible ", res.StatusCode)
		return err
	}

	if err != nil {
		fmt.Println("Error in getting response from api (POST)")
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error in converting response into bytes")
		return err
	}

	var m File
	if err := json.Unmarshal(bodyBytes, &m); err != nil {
		fmt.Println("Error in unparsing json to struct")
		return nil
	}

	for _, mod := range mods {
		mod.ID = m[mod.Hash].ProjectID
		fmt.Println(mod.ID)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(mods))

	for _, mod := range mods {
		wg.Add(1)
		go func(m *Mod) {
			defer wg.Done()
			errChan <- UpdateMod(m, "1.16.5")
		}(mod)

		go func() {
			wg.Wait()
			close(errChan)
		}()

		for err := range errChan {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckInModrinth :- HElPER FUNCTION (To check if a mod is listed in modrinth)
func CheckIfInModrinth(mods map[string]*Mod) {
	for _, mod := range mods {
		url := fmt.Sprintf("%s/version_file/%s", baseURL, mod.Hash)
		res, err := http.Get(url)
		if res.StatusCode == http.StatusNotFound {
			mod.IsModrinth = false
			continue
		}
		mod.IsModrinth = true
		if err != nil {
			fmt.Println("Error in identifying the non-modrinth file")
			return
		}
		res.Body.Close()
	}
}

func UpdateMod(mod *Mod, game_version string) error {
	url := fmt.Sprintf("%s/project/%s/version?loaders=[\"fabric\"]&game_versions=[\"%s\"]", baseURL, mod.ID, game_version)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in getting response from modrinth for Update")
		return err
	}
	defer res.Body.Close()

	fmt.Println(res)
	return nil
}
