package main

import (
	"fmt"
	"log"
	"os"
	"vk-books/pkg/cmd"
	"vk-books/pkg/config"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func main() {

	if err := util.CreateFilesAndFolders(); err != nil {
		fmt.Println("Error creating files/folders:", err)
		os.Exit(1)
	}

	b := db.Books{}

	err := b.ReadFromFile(config.LocalFile)
	if err != nil {
		log.Fatalf("Fatal error: failed to load walkings database: %v", err)
	}

	cmd.CommandLine(&b)
}
