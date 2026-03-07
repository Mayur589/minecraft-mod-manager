package main

import (
	"fmt"
	"log"
	"sync"

	"minecraft-mod-updater/core"
)

func main() {
	const path string = "/Users/mayur/Library/Application Support/minecraft/mods"
	const downloadPath string = "/Users/mayur/Downloads"
	mods, err := core.GiveHash(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := core.GetModIds(mods); err != nil {
		fmt.Println(err)
	}

	// Concurrent Updates with a WaitGroup
	var wg sync.WaitGroup
	errChan := make(chan error, len(mods))

	for _, mod := range mods {
		if !mod.IsModrinth {
			continue
		}

		wg.Add(1)
		go func(m *core.Mod) {
			defer wg.Done()
			if err := core.UpdateMod(m, "1.20.1"); err != nil {
				errChan <- err
			}
		}(mod)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			fmt.Println(fmt.Errorf("error: %w", err))
		}
	}
}
