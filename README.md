# VK-Books

A command-line book tracker written in Go. Add the books you've read, keep your notes and opinions, and get a dashboard with your reading stats, all stored in a plain JSON file with an automatic backup to a USB drive.

```
------------------------
VK-BOOKS
------------------------

Books: 42   Pages: 11830   Re-reads: 5
Average length: 282 pages

Latest
	#38  Hyperion - Dan Simmons
	#39  Snow Crash - Neal Stephenson
	#40  Neuromancer - William Gibson
	#41  Dune Messiah - Frank Herbert
	#42  Foundation - Isaac Asimov        <- newest, in bold green

Books per year
	2026   14  ████████████████████
	2025   11  ████████████████
	2024    6  ████████

Top genres
	sci-fi    15  ███████████████
	history    9  █████████

Top authors
	Frank Herbert    4  ████
	Isaac Asimov     3  ███

Most re-read
	#7  Dune - Frank Herbert  (3 times)
```

## Features

- Add, update and delete books through interactive prompts, with each field pre-filled so you only change what you need
- Input validation: name and author are required, pages and read count must be positive numbers, dates must be `dd.mm.yyyy`
- Language guess from the title (Cyrillic titles suggest Russian)
- Search by title (case-insensitive, partial matches)
- Full history grouped by year, with your opinions underneath
- Dashboard with totals, latest books, books per year, top genres and authors, average length and your most re-read book
- Press Ctrl+C at any prompt to cancel without saving anything
- Atomic writes, so a crash mid-save can't leave a half-written file
- Best-effort backup to an external drive (a failed backup never blocks saving)

## Requirements

- Go 1.21 or newer
- Linux (the backup drive mount is aimed at Linux)

## Installation

```sh
git clone https://github.com/VkHyperNova/VK-Books.git
cd VK-Books
go build -o books .
```

Or run it directly without building:

```sh
go run . <command> [arguments]
```

## Usage

```
books [command] [arguments]
```

| Command                | Description                           |
| ---------------------- | ------------------------------------- |
| *(none)*               | Show the dashboard                    |
| `add`, `a`             | Add a new book                        |
| `update`, `u <id>`     | Update the book with the given id     |
| `delete`, `d <id>`     | Delete the book with the given id     |
| `search`, `s <name>`   | Find books by name                    |
| `history`, `h`         | Show all books, grouped by year       |
| `help`, `-h`, `--help` | Show help                             |

### Examples

```sh
books                  # dashboard
books add              # prompts for name, author, pages, ...
books update 3         # edit book #3, fields pre-filled
books delete 3         # asks for confirmation first
books search dune      # finds "Dune", "Dune Messiah", ...
books history          # everything, grouped by year
```

When adding or updating, each field is shown with a suggested value that you can edit or accept with Enter. Deleting asks `(y/N)` and does nothing unless you answer `y` or `yes`.

## Data

Books are saved as JSON. The path is built from the `baseLocal`, `baseDB` and `file` constants in `main.go` and is **relative to the directory you run the program from**, so run it from the same place each time (or change the constants).

```json
{
  "books": [
    {
      "id": 1,
      "book": "Dune",
      "author": "Frank Herbert",
      "pages": "612",
      "readcount": "3",
      "genre": "sci-fi",
      "language": "English",
      "opinion": "A classic.",
      "date": "14.01.2026"
    }
  ]
}
```

| Field       | Meaning                                                              |
| ----------- | -------------------------------------------------------------------- |
| `id`        | Unique number, assigned automatically (highest id + 1)               |
| `book`      | Title                                                                |
| `author`    | Author                                                               |
| `pages`     | Page count                                                           |
| `readcount` | How many times you've read it (anything above 1 counts as re-reads)  |
| `genre`     | Genre, used for the top genres list                                  |
| `language`  | Language of the book                                                 |
| `opinion`   | Your notes, shown in the history                                     |
| `date`      | Date you read it (`dd.mm.yyyy`), used for books per year             |

### Backup

After every save the program tries to mount the drive labelled `VK DATA` (using [vkmount](https://github.com/VkHyperNova/vkmount)) and copies the data file to `/media/veikko/VK DATA/<baseLocal>/BOOKS/books.json`.

The local save is the one that has to succeed. If the drive isn't plugged in or can't be mounted, you'll see a yellow `Backup failed: ...` notice and your data is still safely saved locally.

To use this on your own machine, change `driveLabel` and `driveMountPoint` in `main.go`.

## Project structure

```
├── main.go           # Entry point, constants, usage text, command dispatch
├── books.go          # Book types, add/update/delete/search, prompts and validation
├── storage.go        # Load, Save, backup, atomic file writes
├── report.go         # Dashboard and history rendering
├── books_test.go     # Tests for the commands and input handling
└── storage_test.go   # Tests for loading, saving and backup
```

## Tests

```sh
go test ./...
```

The tests use a temporary directory, a stubbed drive mount and a scripted stand-in for the interactive prompt, so they never touch your real data, your backup drive or the keyboard.

## License

No license specified. All rights reserved by the author unless stated otherwise.