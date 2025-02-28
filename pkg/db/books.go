package db

import (
	"encoding/json"
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

func (b *Books) ReadFromFile(path string) {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byteValue, b)
	if err != nil {
		panic(err)
	}
}

func (b *Books) Add(newBook Book) error {

	// Append
	b.BOOKS = append(b.BOOKS, newBook)

	// Format The Books JSON
	formattedJSON, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}

	// Save
	err = SaveToFile(config.LocalPath, formattedJSON)
	if err != nil {
		return err
	}

	// Backup Save
	err = SaveToFile(config.BackupPath, formattedJSON)
	if err != nil {
		return err
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

func SaveToFile(path string, newBook []byte) error {

	err := os.WriteFile(path, newBook, 0644)
	if err != nil {
		return err
	}

	return nil
}
