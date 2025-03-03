package cmd

import (
	"fmt"
	"log"
	"strings"
	"vk-books/pkg/config"
	"vk-books/pkg/db"
)

func CommandLine(books *db.Books) {

	var userInput string = ""
	var userInputID int = 0

	fmt.Print("=> ")

	fmt.Scanln(&userInput, &userInputID)

	userInput = strings.ToLower(userInput)

	switch userInput {
	case "a", "add":
		newBook := books.UserInput(config.AddSuggestions)
		err := books.Add(newBook)
		if err != nil {
			log.Fatal(err)
		}

		CommandLine(books)
	case "u", "update":
		index, updateSuggestions := books.FindBook(userInputID)
		if index != -1 {
			updatedBook := books.UserInput(updateSuggestions)
			err := books.Update(index, updatedBook)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("Book %d not found\n", userInputID)
		}

		CommandLine(books)
	case "d", "delete":
		index, _ := books.FindBook(userInputID)
		if index != -1 {
			books.Delete(index)
		}

		CommandLine(books)
	default:
		CommandLine(books)
	}


}
