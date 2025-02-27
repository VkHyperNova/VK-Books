package db

import (
	"encoding/json"
	"io"
	"os"
	"vk-books/pkg/config"
	"vk-books/pkg/util"

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
	
	util.CreateDirectoryIfNotExists("BOOKS")

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

func (b *Books) Insert(newBook Book) error {

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

	bookName := PromptWithSuggestion("Book Name:", "")
	author := PromptWithSuggestion("Author:", "")
	pages := PromptWithSuggestion("Pages:", "")
	readCount := PromptWithSuggestion("Read Count:", "1")
	genre := PromptWithSuggestion("Genre:", "")
	language := PromptWithSuggestion("Language:", "English")
	opinion := PromptWithSuggestion("Opinion:", "")


	return Book{
		ID:        id,
		BOOK:      bookName,
		AUTHOR:    author,
		PAGES:     pages,
		READCOUNT: readCount,
		GENRE:     genre,
		LANGUAGE:  language,
		OPINION:   opinion,
		DATE:      "2022-01-01",
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


