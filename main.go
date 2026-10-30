package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	Bold   = "\033[1m"
	Italic = "\033[3m"

	driveLabel      = "VK DATA"
	driveMountPoint = "/media/veikko/VK DATA"
	file            = "books.json"
	baseDB          = "BOOKS"
	baseLocal       = "DATABASES"
	// baseLocal = "DEVELOPMENT"

	usage = Bold + "Usage:" + Reset + `
  books [command] [arguments]

` + Bold + "Commands:" + Reset + `
  add, a             Add a new book
  update, u <id>     Update the book with the given id
  delete, d <id>     Delete the book with the given id
  search, s <name>   Find books by name (case-insensitive)
  history, h         Show all books, grouped by year
  help, -h, --help   Show this help

With no command, the dashboard is shown.

` + Bold + "Examples:" + Reset + `
  books add
  books update 3
  books delete 3
  books search dune
  books history
`
)

var (
	LocalFile  = filepath.Join(baseLocal, baseDB, file)
	BackupFile = filepath.Join(driveMountPoint, baseLocal, baseDB, file)
)

var notices []string

func notice(msg string) { notices = append(notices, msg) }

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			fmt.Print(usage)
			return
		}
	}

	db := Books{}
	if err := db.Load(LocalFile); err != nil {
		log.Fatalf("load from file: %v", err)
	}

	// What gets printed after the command runs. Default is the dashboard;
	// commands can swap it or set it to nil to print nothing extra.
	show := db.Report

	var err error
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "a", "add":
			err = db.Add()
		case "u", "update":
			err = db.Update()
		case "d", "delete":
			err = db.Delete()
		case "h", "history":
			show = db.PrintHistory
		case "s", "search":
			err = db.Search(strings.Join(os.Args[2:], " "))
			show = nil // Search prints its own results
		default:
			fmt.Printf("%serror: unknown command: %s%s\n\n", Red, os.Args[1], Reset)
			fmt.Print(usage)
			os.Exit(1)
		}
	}

	// os.Exit skips defers, so restore the terminal explicitly
	closePrompt()

	// Ctrl+C / Ctrl+D at a prompt is not an error
	if errors.Is(err, errCancelled) {
		fmt.Println(Gray + "Cancelled." + Reset)
		err = nil
	}

	if show != nil {
		show()
	}

	for _, n := range notices {
		fmt.Println(Yellow + n + Reset)
	}

	if err != nil {
		fmt.Printf("%serror: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}
}
