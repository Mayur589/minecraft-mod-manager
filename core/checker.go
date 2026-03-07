package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://api.modrinth.com/v2"

func GetModIds(mods map[string]*Mod) error {
	hashes := make([]string, 0, len(mods))
	for _, mod := range mods {
		hashes = append(hashes, mod.Hash)
	}

	payload, _ := json.Marshal(map[string]any{
		"hashes":    hashes,
		"algorithm": "sha1",
	})

	resp, err := http.Post(baseURL+"/version_files", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("bulk check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var results map[string]struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return err
	}

	// Update local mod state based on results
	for _, mod := range mods {
		if data, found := results[mod.Hash]; found {
			mod.IsModrinth = true
			mod.ID = data.ProjectID
		} else {
			mod.IsModrinth = false
		}
	}

	return nil
}

// UpdateMod :- Update the mod to preffered version
func UpdateMod(mod *Mod, gameVersion string) error {
	if !mod.IsModrinth {
		fmt.Println(mod.Name, " is not in the modrinth")
		return nil
	}

	params := url.Values{}
	params.Add("loaders", "[\"fabric\"]")
	params.Add("game_versions", fmt.Sprintf("[\"%s\"]", gameVersion))

	url := fmt.Sprintf("%s/project/%s/version?%s", baseURL, mod.ID, params.Encode())
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in getting response from modrinth for Update")
		return err
	}
	defer res.Body.Close()

	var results []struct {
		Files []struct {
			URL      string `json:"url"`
			FileName string `json:"filename"`
		} `json:"files"`
	}

	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		fmt.Println(res.StatusCode)
		return fmt.Errorf("failed to decode JSON for %s: %w", mod.Name, err)
	}
	// Safety check: Did we actually get any versions back?
	if len(results) > 0 && len(results[0].Files) > 0 {
		mod.DownloadURL = results[0].Files[0].URL
	} else {
		fmt.Printf("No compatible versions found for %s on %s\n", mod.Name, gameVersion)
	}

	return nil
}
