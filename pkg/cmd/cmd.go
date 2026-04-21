package cmd

import (
	"fmt"
	"os"
	"strings"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func CommandLine(b *db.Books) {
	for {

		b.PrintCLI()

		var cmd string
		var id int

		fmt.Print("=> ")

		fmt.Scanln(&cmd, &id)

		cmd = strings.ToLower(cmd)

		switch cmd {
		case "a", "add":
			err := b.Add()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Book added!")
			}

		case "u", "update":
			err := b.Update(id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%d is updated!\n", id)
			}
		case "d", "delete":
			err := b.Delete(id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%d is deleted!\n", id)
			}
		case "showall":
			b.PrintAllBooks()
		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)
		case "":
			fmt.Println("Please enter a command.")
		default:
			b.Search(cmd)
		}
	}
}
