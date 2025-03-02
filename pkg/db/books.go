package db

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
	"vk-books/pkg/config"
	"vk-books/pkg/util"
)

type Book struct {
	ID        int    `json:"id"`
	BOOK      string `json:"book"`
	AUTHOR    string `json:"author"`
	PAGES     string `json:"pages"`
	READCOUNT string `json:"readcount"`
	GENRE     string `json:"genre"`
	LANGUAGE  string `json:"language"`
	OPINION   string `json:"opinion"`
	DATE      string `json:"date"`
}

type Books struct {
	BOOKS []Book `json:"books"`
}

func (b *Books) ReadFromFile(path string) error {

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer file.Close()

	// Read entire file content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(byteValue, b); err != nil {
		return fmt.Errorf("error parsing JSON from file %s: %w", path, err)
	}

	return nil
}

func (b *Books) Add(newBook Book) error {
	b.BOOKS = append(b.BOOKS, newBook)
	return b.Save()
}

func (b *Books) Save() error {

	// Format JSON
	books, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}

	// Save
	err = os.WriteFile(config.LocalPath, books, 0644)
	if err != nil {
		return err
	}

	// Save Backup
	err = os.WriteFile(config.BackupPath, books, 0644)
	if err != nil {
		return err
	}

	// Save Backup with Date
	err = os.WriteFile(config.BackupPathWithDate, books, 0644)
	if err != nil {
		return err
	}

	return nil
}

func UserInput(id int) Book {

	book := util.PromptWithSuggestion("Book Name:", "")
	author := util.PromptWithSuggestion("Author:", "")
	pages := util.PromptWithSuggestion("Pages:", "")
	readCount := util.PromptWithSuggestion("Read Count:", "1")
	genre := util.PromptWithSuggestion("Genre:", "")
	language := util.PromptWithSuggestion("Language:", "English")
	opinion := util.PromptWithSuggestion("Opinion:", "")
	date := util.PromptWithSuggestion("Date:", time.Now().Format("02.01.2006"))

	return Book{
		ID:        id,
		BOOK:      book,
		AUTHOR:    author,
		PAGES:     pages,
		READCOUNT: readCount,
		GENRE:     genre,
		LANGUAGE:  language,
		OPINION:   opinion,
		DATE:      date,
	}
}

func (b *Books) NewID() int {

	maxID := 0

	for _, book := range b.BOOKS {
		if book.ID > maxID {
			maxID = book.ID
		}
	}

	return maxID + 1
}

func HandleBackupRestore(books *Books) error {
	useBackup := flag.Bool("backup", false, "Use backup database file")
	flag.Parse()

	if *useBackup {
		fmt.Println("Using backup database file.")

		// Read from backup
		if err := books.ReadFromFile(config.BackupPath); err != nil {
			return fmt.Errorf("failed to load books database from backup: %w", err)
		}

		// Format JSON
		formattedBooks, err := json.MarshalIndent(books, "", "  ")
		if err != nil {
			return err
		}

		// Save the backup database as the main database
		if err := os.WriteFile(config.LocalPath, formattedBooks, 0644); err != nil {
			return fmt.Errorf("failed to save JSON file from BACKUP database: %w", err)
		}
	}
	return nil
}
