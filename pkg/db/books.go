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
	Id        int    `json:"id"`
	Name      string `json:"book"`
	Author    string `json:"author"`
	Pages     string `json:"pages"`
	ReadCount string `json:"readcount"`
	Genre     string `json:"genre"`
	Language  string `json:"language"`
	Opinion   string `json:"opinion"`
	Date      string `json:"date"`
}

type Books struct {
	Books []Book `json:"books"`
}

/* Operations */

func (b *Books) LoadFromFile(path string) error {

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

func (b *Books) PrintDashboard() {

	// Program information
	fmt.Println(color.Cyan + color.Bold + "------------------------" + color.Reset)
	fmt.Println(color.Cyan + color.Bold + "VK-BOOKS" + color.Reset)
	fmt.Println(color.Cyan + color.Bold + "------------------------" + color.Reset)
	

	// Print total pages and book count
	b.PrintLatest(3)
	b.PrintSummary()
}

func (b *Books) Add() error {

	newBook, err := b.promptBookInput(Book{})
	if err != nil {
		return err
	}

	// Add unique ID
	newBook.Id = b.nextID()

	// Add
	b.Books = append(b.Books, newBook)

	return b.saveToFile()
}

func (b *Books) Update(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := b.indexOf(id)
	if err != nil {
		return err
	}

	updated, err := b.promptBookInput((b.Books)[index])
	if err != nil {
		return err
	}

	(b.Books)[index] = updated

	return b.saveToFile()
}

func (b *Books) Delete(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := b.indexOf(id)
	if err != nil {
		return err
	}

	confirm := util.PromptConfirm()
	if !confirm {
		return fmt.Errorf("Abort")
	}

	b.Books = append((b.Books)[:index], (b.Books)[index+1:]...)

	return b.saveToFile()
}

func (b *Books) PrintLatest(numberOfBooks int) {

	fmt.Printf("%s\n", color.Yellow+color.Bold+color.Italic+"Latest Books: "+color.Reset)

	if numberOfBooks <= 0 {
		return
	}
	start := len(b.Books) - numberOfBooks
	if start < 0 {
		start = 0
	}
	for _, book := range b.Books[start:] {
		fmt.Println(b.formatBook(book))
	}
}

func (b *Books) PrintSummary() {

	totalPagesRead := fmt.Sprintf("\n%d Pages | %d Books\n", b.totalPages(), len(b.Books))
	fmt.Print(color.Yellow + totalPagesRead + color.Reset)
}

func (b *Books) Search(searchBook string) {
	fmt.Printf("%s\n", color.Yellow+color.Bold+color.Italic+"Found Books: "+color.Reset)
	for _, book := range b.Books {
		if strings.Contains(strings.ToLower(book.Name), strings.ToLower(searchBook)) {
			fmt.Println(b.formatBook(book))
		}
	}
}

func (b *Books) PrintAll() {
	for _, book := range b.Books {
		fmt.Println(b.formatBook(book))
	}
}

/* Helpers */

func (b *Books) totalPages() int {
	totalPages := 0
	for _, book := range b.Books {
		pages, err := strconv.Atoi(book.Pages)
		if err != nil {
			fmt.Print(book)
			fmt.Println(err)
		}
		totalPages = totalPages + pages
	}

	return totalPages
}

func (b *Books) formatBook(book Book) string {
	bookID := fmt.Sprintf("%s%v%s", color.Yellow, book.Id, color.Reset)
	bookName := fmt.Sprintf("%s\"%s\"%s", color.Green, book.Name, color.Reset)
	bookAuthor := fmt.Sprintf("%s by %s%s", color.Cyan, book.Author, color.Reset)
	bookPages := fmt.Sprintf("(%s pages)", book.Pages)
	bookReadCount := fmt.Sprintf("[%s]", book.ReadCount)
	bookGenre := fmt.Sprintf("(%s)", book.Genre)
	bookLanguage := fmt.Sprintf("(%s)", book.Language)
	bookOpinion := fmt.Sprintf("(%s)", book.Opinion)

	return fmt.Sprintf("%s %s %s %s\t%s%s %s %s %s%s %s",
		bookID, bookName, bookAuthor,
		color.Purple, bookPages, bookReadCount,
		bookGenre, bookLanguage, bookOpinion, color.Reset,
		book.Date)
}

func (b *Books) saveToFile() error {

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

func (b *Books) indexOf(id int) (int, error) {
	for i, books := range b.Books {
		if books.Id == id {
			fmt.Println(books)
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

func (b *Books) promptBookInput(suggestions Book) (Book, error) {
	prompts := []struct {
		label  string
		target *string
	}{
		{"Book Name:", &suggestions.Name},
		{"Author:", &suggestions.Author},
		{"Pages Count:", &suggestions.Pages},
		{"Read Count:", &suggestions.ReadCount},
		{"Genre:", &suggestions.Genre},
		{"Language:", &suggestions.Language},
		{"Opinion:", &suggestions.Opinion},
		{"Date:", &suggestions.Date},
	}

	for _, p := range prompts {
		val, err := util.InputWithSuggestion(p.label, *p.target)
		if err != nil {
			return Book{}, err
		}
		*p.target = val
	}

	if suggestions.Language == "" {
		suggestions.Language = util.DetectLanguage(suggestions.Name)
	}
	if suggestions.Date == "" {
		suggestions.Date = time.Now().Format("02.01.2006")
	}

	return suggestions, nil
}

func (b *Books) nextID() int {

	maxID := 0

	for _, book := range b.Books {
		if book.Id > maxID {
			maxID = book.Id
		}
	}

	return maxID + 1
}
