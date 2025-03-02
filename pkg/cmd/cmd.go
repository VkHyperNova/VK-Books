package cmd

import (
	"fmt"
	"log"
	"strings"
	"vk-books/pkg/db"
)

func CommandLine(books *db.Books) {

	var userInput string = ""
	var userInputID int = 0

	fmt.Print("=> ")

	fmt.Scanln(&userInput, &userInputID)

	userInput = strings.ToLower(userInput)

	switch userInput {
	case "a", "add", "insert":
		id := books.NewID()
		newBook := db.UserInput(id)
		err := books.Add(newBook)
		if err != nil {
			log.Fatal(err)
		}
	case "u", "update":
		// Update
	case "d", "delete":
		// Delete
	}

}
