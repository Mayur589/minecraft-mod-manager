package main

import (
	"fmt"
	"log"
	"minecraft-mod-updater/core"
	"sync"
)

func main() {
	fmt.Println("Hello from main")
	const path string = "/Users/mayur/Library/Application Support/minecraft/mods"
	mods, err := core.GiveHash(path)

	if err != nil {

		fmt.Println("Err in readinf folder")
		log.Fatal(err)

	}

	// fmt.Println(mods)

	for _, modName := range mods {
		// fmt.Printf("%T\n", modName.Hash)
		var wg sync.WaitGroup
		fmt.Printf("%s: \n", modName.ID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			core.GetUpdate(modName.Hash)
		}()

		wg.Wait()
	}
}
