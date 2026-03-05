package main

import (
	"log"

	"minecraft-mod-updater/core"
)

func main() {
	const path string = "/Users/mayur/Library/Application Support/minecraft/mods"
	const downloadPath string = "/Users/mayur/Downloads"
	mods, err := core.GiveHash(path)
	if err != nil {
		log.Fatal(err)
	}

	core.GetModIds(mods)

	// for _, mod := range mods {
	// 	fmt.Println(mod.Hash)
	// }

	// if err = core.CheckInModrinth(mods, "1.16.5", downloadPath); err != nil {
	// 	fmt.Println("Error in downloading update")
	// 	log.Fatal(err)
	// }
}
