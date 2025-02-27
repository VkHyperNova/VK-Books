package main

import (
	"vk-books/pkg/cmd"
	"vk-books/pkg/config"
	"vk-books/pkg/db"
)

func main() {

	// Load Database
    books := db.Books{}
	books.ReadFromFile(config.LocalPath)

	// Start
	cmd.CommandLine(&books)	
}
