package cmd

import (
	"fmt"
	"os"
	"strings"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func Run(b *db.Books) {
	for {

		b.PrintDashboard()

		var command string
		var id int

		fmt.Print("=> ")

		fmt.Scanln(&command, &id)

		command = strings.ToLower(command)

		switch command {
		case "a", "add":
			if err := b.Add(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Added!")
			}
		case "u", "update":
			if err := b.Update(id); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Updated!")
			}
		case "d", "delete":
			if err := b.Delete(id); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Deleted!")
			}
		case "showall", "all":
			b.PrintAll()
		case "q", "quit":
			util.ClearTerminal()
			os.Exit(0)
		case "":
			fmt.Println("(add, update, delete, all, q or type book name)")
		default:
			b.Search(command)
		}
	}
}
