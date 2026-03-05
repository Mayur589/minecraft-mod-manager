// loadMods - used to load all the mods in mods
// Mods structure - map[string(name of mod)](ptr)Mod (model.go)

package core

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GiveHash scans a directory for .jar files and returns a map of mod name -> Mod with hash.
func GiveHash(dir string) (map[string]*Mod, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %q: %w", dir, err)
	}

	mods := make(map[string]*Mod)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jar") {
			continue
		}

		modPath := filepath.Join(dir, entry.Name())
		if err := processFile(modPath, mods); err != nil {
			return nil, fmt.Errorf("processing mod %q: %w", entry.Name(), err)
		}
	}

	return mods, nil
}

// processFile opens a .jar file, parses its metadata, and adds it to the mods map.
func processFile(modPath string, mods map[string]*Mod) error {
	r, err := zip.OpenReader(modPath)
	if err != nil {
		return fmt.Errorf("opening jar: %w", err)
	}
	defer r.Close()

	mod, err := parseFabricMod(r, modPath)
	if err != nil {
		return err
	}
	if mod == nil {
		return nil // not a fabric mod, skip
	}

	hash, err := getHash(modPath)
	if err != nil {
		return fmt.Errorf("hashing mod: %w", err)
	}

	mod.Hash = hash
	mods[mod.Name] = mod
	return nil
}

// parseFabricMod finds and parses fabric.mod.json inside a jar, returning a Mod or nil if not found.
func parseFabricMod(r *zip.ReadCloser, modPath string) (*Mod, error) {
	for _, f := range r.File {
		if !strings.Contains(f.Name, "fabric.mod.json") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("opening fabric.mod.json: %w", err)
		}
		defer rc.Close()

		content, err := io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("reading fabric.mod.json: %w", err)
		}

		var meta jsonModFile
		if err := json.Unmarshal(content, &meta); err != nil {
			return nil, fmt.Errorf("parsing fabric.mod.json: %w", err)
		}

		return &Mod{
			Name:    meta.Name,
			ID:      meta.ID,
			Version: meta.Version,
			Path:    modPath,
		}, nil
	}

	return nil, nil
}

func getHash(modPath string) (string, error) {
	f, err := os.Open(modPath)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashing file: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
