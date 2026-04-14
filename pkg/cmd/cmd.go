package cmd

import (
	"fmt"
	"os"
	"strings"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func CommandLine(books *db.Books) {
	for {

		books.PrintCLI()

		var cmd string
		var id int

		fmt.Print("=> ")

		fmt.Scanln(&cmd, &id)

		cmd = strings.ToLower(cmd)

		switch cmd {
		case "a", "add":
			err := books.Add()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Book added!")
			}

		case "u", "update":
			err := books.Update(id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%d is updated!\n", id)
			}
		case "d", "delete":
			err := books.Delete(id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%d is deleted!\n", id)
			}
		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
		}
	}
}
