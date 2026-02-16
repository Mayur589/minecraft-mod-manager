package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetUpdate(hash string) (bool, error) {
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

	// fmt.Println(m.GameVersions, m.Name, m.VersionNumber)
	fmt.Println(m.Files)

	// fmt.Printf("%s", body)
	return true, nil
}
