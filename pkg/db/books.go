package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"vk-books/pkg/config"

	"github.com/peterh/liner"
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

	// Append
	b.BOOKS = append(b.BOOKS, newBook)

	// Save
	err := b.SaveToFile(config.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to save JSON file to LocalPath: %w", err)
	}

	// Backup Save
	err = b.SaveToFile(config.BackupPath)
	if err != nil {
		return fmt.Errorf("failed to save JSON file to BackupPath: %w", err)
	}

	// Success
	return nil
}

func GetUserInput(id int) Book {

	book := PromptWithSuggestion("Book Name:", "")
	author := PromptWithSuggestion("Author:", "")
	pages := PromptWithSuggestion("Pages:", "")
	readCount := PromptWithSuggestion("Read Count:", "1")
	genre := PromptWithSuggestion("Genre:", "")
	language := PromptWithSuggestion("Language:", "English")
	opinion := PromptWithSuggestion("Opinion:", "")
	date := PromptWithSuggestion("Date:", time.Now().Format("02.01.2006"))

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

func PromptWithSuggestion(name string, suggestion string) string {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", suggestion, -1)
	if err != nil {
		panic(err)
	}

	return input
}

func (b *Books) GenerateUniqueID() int {

	maxID := 0

	for _, book := range b.BOOKS {
		if book.ID > maxID {
			maxID = book.ID
		}
	}

	return maxID + 1
}

func (b *Books) SaveToFile(path string) error {

	newBook, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, newBook, 0644)
	if err != nil {
		return err
	}

	return nil
}
