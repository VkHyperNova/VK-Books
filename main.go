package main

import (
	"flag"
	"fmt"
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
	if err := HandleBackupRestore(&books); err != nil {
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

// HandleBackupRestore checks if the backup flag is set and restores the backup database if needed.
func HandleBackupRestore(books *db.Books) error {
	useBackup := flag.Bool("backup", false, "Use backup database file")
	flag.Parse()

	if *useBackup {
		fmt.Println("Using backup database file.")

		// Read from backup
		if err := books.ReadFromFile(config.BackupPath); err != nil {
			return fmt.Errorf("failed to load books database from backup: %w", err)
		}

		// Save the backup database as the main database
		if err := books.SaveToFile(config.LocalPath); err != nil {
			return fmt.Errorf("failed to save JSON file from BACKUP database: %w", err)
		}
	}
	return nil
}
