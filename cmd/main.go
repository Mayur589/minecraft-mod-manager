package main

import (
	"fmt"
	"log"
	"minecraft-mod-updater/core"
)

func main() {
	const path string = "/Users/mayur/Library/Application Support/minecraft/mods"
	mods, err := core.GiveHash(path)

	if err != nil {

		fmt.Println("Err in readinf folder")
		log.Fatal(err)

	}
	core.CheckFunInModrinth(mods)

	for _, mod := range *mods {
		fmt.Printf("%v, %v\n", mod.ID, mod.IsModrinth)
	}
}
