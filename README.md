# VK-Books

A command-line book tracker written in Go. Manage your book collection or reading list through an interactive terminal shell, with data persisted locally.

## Tech Stack

- [Go](https://go.dev/) 1.21
- [peterh/liner](https://github.com/peterh/liner) — interactive line editing / readline-style shell
- [mattn/go-runewidth](https://github.com/mattn/go-runewidth) — terminal width handling

## Project Structure

```
├── pkg/
│   ├── cmd/       # Command handling / interactive shell logic
│   ├── config/    # Configuration (e.g. local storage file path)
│   ├── db/        # Books data model and load/save logic
│   └── util/      # Utility helpers (local storage init, line reader)
├── main.go        # Entry point
├── go.mod
└── go.sum
```

## How It Works

On startup, `VK-Books`:

1. Initializes local storage (`util.InitLocalStorage`).
2. Loads existing book data from a local file (`db.Books.LoadFromFile`).
3. Starts an interactive command shell (`cmd.Run`) where you can add, view, and manage books.
4. Cleans up the interactive line reader on exit.

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

```bash
git clone https://github.com/VkHyperNova/VK-Books.git
cd VK-Books
go build -o vk-books .
```

### Usage

Run the built binary:

```bash
./vk-books
```

Or run directly with Go:

```bash
go run main.go
```

This launches an interactive shell where you can enter commands to manage your book collection, which is persisted locally between sessions.

## License

No license specified — all rights reserved by the author unless stated otherwise.
