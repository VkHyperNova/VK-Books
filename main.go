package main

import (
	"log"
	"vk-books/pkg/cmd"
	"vk-books/pkg/config"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func main() {

	// Create necessary files
	err := util.CreateNecessaryFiles()
	if err != nil {
		log.Fatalf("Fatal error: failed to create necessary files: %v", err)
	}

	// Initialize Books database
	books := db.Books{}

	// Handle backup restoration
	if err := db.HandleBackupRestore(&books); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	// Reload Database
	err = books.ReadFromFile(config.LocalPath)
	if err != nil {
		log.Fatalf("Fatal error: failed to load books database: %v", err)
	}

	// Start
	cmd.CommandLine(&books)
}