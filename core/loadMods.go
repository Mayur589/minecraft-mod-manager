package core

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Mod struct {
	ID      string
	Version string
	Path    string
	Hash    string
}

type jsonModFile struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

func GiveHash(dir string) (map[string]*Mod, error) {
	path := dir

	// map
	mods := make(map[string]*Mod)

	folder, err := os.ReadDir(path)

	if err != nil {

		return nil, err
	}

	for _, m := range folder {
		if !m.IsDir() && strings.HasSuffix(m.Name(), ".jar") {
			modPath := path + "/" + m.Name()

			r, err := zip.OpenReader(modPath)

			if err != nil {
				log.Fatal(err)
			}

			modName := readFile(r, mods, &modPath)
			mods[modName].Hash = getHash(modPath)
			r.Close()

		}
	}
	return mods, nil

}

func readFile(r *zip.ReadCloser, modMap map[string]*Mod, modPath *string) string {
	for _, f := range r.File {
		if strings.Contains(f.Name, "fabric.mod.json") {
			fileReader, err := f.Open()

			if err != nil {
				fmt.Println("Error in reading file")
				log.Fatal(err)
			}
			defer fileReader.Close()

			content, err := io.ReadAll(fileReader)
			if err != nil {
				fmt.Println("Error in reading content of file")
				log.Fatal(err)
			}

			var m jsonModFile
			e := json.Unmarshal(content, &m)

			if e != nil {
				fmt.Println("Error in converting to json")
				log.Fatal(e)
			}

			modMap[m.Name] = &Mod{
				ID:      m.ID,
				Version: m.Version,
				Path:    *modPath,
			}
			return m.Name
		}
	}
	return ""
}

func getHash(modPath string) string {
	file, err := os.Open(modPath)

	if err != nil {

		log.Fatal(err)
	}
	defer file.Close()

	hasher := sha1.New()

	if _, err := io.Copy(hasher, file); err != nil {

		log.Fatal(err)
	}

	hashInBytes := hasher.Sum(nil)

	hashInHex := hex.EncodeToString(hashInBytes)
	return hashInHex
}
