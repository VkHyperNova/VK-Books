package cmd

import (
	"fmt"
	"vk-books/pkg/color"
	"vk-books/pkg/db"
	"vk-books/pkg/util"
)

func Run(b *db.Books) {

	b.PrintDashboard()

	for {

		fmt.Println(color.Yellow + "\n< add, update, delete, import, export, unmount, stats, history, q >" + color.Reset)

		command, id := util.ReadInput()
		if command == "" {
			fmt.Println("Enter Pressed!")
			continue
		}

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
		case "import", "i":
			if err := b.Import(); err != nil {
				fmt.Println(err)
			}
		case "export", "e":
			if err := b.Export(); err != nil {
				fmt.Println(err)
			}
		case "unmount":
			if err := util.UnmountDrive(); err != nil {
				fmt.Println(err)
			}
		case "stats":
			b.Stats()
		case "history", "h":
			b.History()
		case "q", "quit":
			return
		case "search", "s":
			b.Search()
		}
	}
}
