package cmd

import (
	"fmt"
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
		id := books.GenerateUniqueID()
		newBook := db.GetUserInput(id)
		books.Add(newBook)
	case "u", "update":
		// Update
	case "d", "delete":
		// Delete
	}

}
