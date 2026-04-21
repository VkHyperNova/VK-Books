package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"vk-books/pkg/color"
	"vk-books/pkg/config"
	"vk-books/pkg/util"
)

type Book struct {
	ID        int    `json:"id"`
	NAME      string `json:"book"`
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

func (b *Books) PrintCLI() {

	// Program information
	fmt.Println(color.Cyan + color.Bold + "------------------------" + color.Reset)
	fmt.Println(color.Cyan + color.Bold + "VK-BOOKS" + color.Reset)
	fmt.Println(color.Cyan + color.Bold + "------------------------" + color.Reset)

	// Print total pages and book count
	b.PrintLatestBooks(3)
	b.PrintStats()
}

func (b *Books) PrintLatestBooks(numberOfBooks int) {
    if numberOfBooks <= 0 {
        return
    }
    start := len(b.BOOKS) - numberOfBooks
    if start < 0 {
        start = 0
    }
    for _, book := range b.BOOKS[start:] {
        fmt.Println(b.formatBook(book))
    }
}

func (b *Books) PrintStats() {

	totalPagesRead := fmt.Sprintf("\n%d Pages | %d Books\n", b.CountPages(), len(b.BOOKS))
	fmt.Print(color.Yellow + totalPagesRead + color.Reset)
}

func (b *Books) CountPages() int {
	totalPages := 0
	for _, book := range b.BOOKS {
		pages, err := strconv.Atoi(book.PAGES)
		if err != nil {
			fmt.Print(book)
			fmt.Println(err)
		}
		totalPages = totalPages + pages
	}

	return totalPages
}

func (b *Books) formatBook(book Book) string {
    bookID := fmt.Sprintf("%s%v%s", color.Yellow, book.ID, color.Reset)
    bookName := fmt.Sprintf("%s\"%s\"%s", color.Green, book.NAME, color.Reset)
    bookAuthor := fmt.Sprintf("%s by %s%s", color.Cyan, book.AUTHOR, color.Reset)
    bookPages := fmt.Sprintf("(%s pages)", book.PAGES)
    bookReadCount := fmt.Sprintf("[%s]", book.READCOUNT)
    bookGenre := fmt.Sprintf("(%s)", book.GENRE)
    bookLanguage := fmt.Sprintf("(%s)", book.LANGUAGE)
    bookOpinion := fmt.Sprintf("(%s)", book.OPINION)

    return fmt.Sprintf("%s %s %s %s\t%s%s %s %s %s%s %s",
        bookID, bookName, bookAuthor,
        color.Purple, bookPages, bookReadCount,
        bookGenre, bookLanguage, bookOpinion, color.Reset,
        book.DATE)
}

func (b *Books) Search(searchBook string) {
    for _, book := range b.BOOKS {
        if strings.Contains(strings.ToLower(book.NAME), strings.ToLower(searchBook)) {
            fmt.Println(b.formatBook(book))
        }
    }
}

func (b *Books) PrintAllBooks() {
    for _, book := range b.BOOKS {
        fmt.Println(b.formatBook(book))
    }
}

func (b *Books) ReadFromFile(path string) error {

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer file.Close()

	// Read entire file contents
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

func (b *Books) Add() error {

	newBook, err := b.GetUserInput(Book{})
	if err != nil {
		return err
	}

	// Add unique ID
	newBook.ID = b.NewID()

	// Add
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
	err = os.WriteFile(config.LocalFile, books, 0644)
	if err != nil {
		return err
	}

	// Save Backup
	err = os.WriteFile(config.BackupFile, books, 0644)
	if err != nil {
		return err
	}

	// Save Backup with Date
	err = os.WriteFile(config.BackupFileWithDate, books, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (b *Books) Update(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := b.findIndex(id)
	if err != nil {
		return err
	}

	updated, err := b.GetUserInput((b.BOOKS)[index])
	if err != nil {
		return err
	}

	(b.BOOKS)[index] = updated

	return b.Save()
}

func (b *Books) Delete(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := b.findIndex(id)
	if err != nil {
		return err
	}

	confirm := util.Confirm()
	if !confirm {
		return fmt.Errorf("Abort")
	}

	b.BOOKS = append((b.BOOKS)[:index], (b.BOOKS)[index+1:]...)

	return b.Save()
}

func (b *Books) findIndex(id int) (int, error) {
	for i, books := range b.BOOKS {
		if books.ID == id {
			fmt.Println(books)
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

func (b *Books) GetUserInput(suggestions Book) (Book, error) {
	prompts := []struct {
		label  string
		target *string
	}{
		{"Book Name:", &suggestions.NAME},
		{"Author:", &suggestions.AUTHOR},
		{"Pages Count:", &suggestions.PAGES},
		{"Read Count:", &suggestions.READCOUNT},
		{"Genre:", &suggestions.GENRE},
		{"Language:", &suggestions.LANGUAGE},
		{"Opinion:", &suggestions.OPINION},
		{"Date:", &suggestions.DATE},
	}

	for _, p := range prompts {
		val, err := util.PromptWithSuggestion(p.label, *p.target)
		if err != nil {
			return Book{}, err
		}
		*p.target = val
	}

	if suggestions.LANGUAGE == "" {
		suggestions.LANGUAGE = util.AutoDetectLanguage(suggestions.NAME)
	}
	if suggestions.DATE == "" {
		suggestions.DATE = time.Now().Format("02.01.2006")
	}

	return suggestions, nil
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
