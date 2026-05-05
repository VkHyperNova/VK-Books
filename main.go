package main

import (
	"log"
	"vk-books/pkg/cmd"
	"vk-books/pkg/config"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func main() {

	defer util.Line.Close()

	if err := util.InitLocalStorage(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	b := db.Books{}

	err := b.LoadFromFile(config.LocalFile)
	if err != nil {
		log.Fatalf("Failed to load books: %v", err)
	}

	cmd.Run(&b)
}
