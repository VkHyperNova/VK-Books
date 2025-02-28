package main

import (
	"vk-books/pkg/cmd"
	"vk-books/pkg/config"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func main() {
	util.CreateDirectoryIfNotExists("BOOKS")

	// Load Database
    books := db.Books{}
	books.ReadFromFile(config.LocalPath)

	// Start
	cmd.CommandLine(&books)	
}
