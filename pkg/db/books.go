package db

import (
	"encoding/json"
	"os"
)

type Book struct {
	ID        int    `json:"id"`
	BOOK      string `json:"book"`
	AUTHOR    string `json:"author"`
	PAGES     int    `json:"pages"`
	READCOUNT int    `json:"readcount"`
	TYPE      string `json:"type"`
	LANGUAGE  string `json:"language"`
	DATE      string `json:"date"`
}

type Books struct {
	BOOKS []Book `json:"books"`
}

func (b *Books) Insert() error {

	id := b.UniqueID()

	newBook := Book{
		ID:        id,
		BOOK:      "Book1",
		AUTHOR:    "Author1",
		PAGES:     100,
		READCOUNT: 0,
		TYPE:      "Fiction",
		LANGUAGE:  "English",
		DATE:      "2022-01-01",
	}

	b.BOOKS = append(b.BOOKS, newBook)

	formattedBooks, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}

	SaveToFile("./BOOKS/books.json", formattedBooks)
	SaveToFile("/media/veikko/VK DATA/DATABASES/BOOKS/books.json", formattedBooks) // Backup

	return nil
}

func (b *Books) UniqueID() int {

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
