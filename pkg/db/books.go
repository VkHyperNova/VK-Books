package db

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
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
	fmt.Println(color.Cyan + "VK-BOOKS" + color.Reset)
	fmt.Println(color.Cyan + "------------------------" + color.Reset)

	// Print all books by genre
	b.PrintBooks()

	// Print total pages and book count
	b.PrintTotalPagesAndBookCount()
}

func (b *Books) PrintTotalPagesAndBookCount() {g
	totalPages := 0
	for _, book := range b.BOOKS {
		pages, err := strconv.Atoi(book.PAGES)
		if err != nil {
			fmt.Print(book)
			fmt.Println(err)
		}
		totalPages = totalPages + pages
	}

	totalPagesRead := fmt.Sprintf("%d Pages | %d Books\n", totalPages, len(b.BOOKS))
	fmt.Print(color.Yellow + totalPagesRead + color.Reset)
}

func (b *Books) PrintBooks() {
	for _, book := range b.BOOKS {
		bookID := fmt.Sprint(color.Yellow, book.ID, color.Reset)
		bookName := fmt.Sprint(color.Green + "\"" + book.NAME + "\"" + color.Reset)
		bookAuthor := fmt.Sprint(color.Cyan + " by " + book.AUTHOR + color.Reset)
		bookPages := fmt.Sprint("(" + book.PAGES +" pages)")
		bookReadCount := fmt.Sprint("[" + book.READCOUNT +"]")
		bookGenre := fmt.Sprint("(" + book.GENRE + ")")
		bookLanguage := fmt.Sprint("(" + book.LANGUAGE + ")")
		bookOpinion := fmt.Sprint("(" + book.OPINION + ")")
		bookDate := fmt.Sprint(book.DATE)
		fmt.Println(bookID, bookName, bookAuthor,color.Purple+ "\t"+bookPages + bookReadCount, bookGenre, bookLanguage, bookOpinion, color.Reset, bookDate)
	}
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

func (b *Books) Update(index int, updatedBook Book) error {

	// Set correct ID
	updatedBook.ID = b.BOOKS[index].ID

	// Update
	b.BOOKS[index] = updatedBook

	return b.Save()
}

func (b *Books) Delete(index int) error {
	b.BOOKS = append((b.BOOKS)[:index], (b.BOOKS)[index+1:]...)
	return b.Save()
}

func (b *Books) FindBook(searchBookID int) (int, []string) {

	for index, foundBook := range b.BOOKS {
		if foundBook.ID == searchBookID {
			return index, []string{foundBook.NAME, foundBook.AUTHOR, foundBook.PAGES, foundBook.READCOUNT, foundBook.GENRE, foundBook.LANGUAGE, foundBook.OPINION, foundBook.DATE}
		}
	}

	return -1, nil
}

func (b *Books) UserInput(suggestions []string) Book {

	var answers []string
	for index, question := range config.Questions {
		input := util.PromptWithSuggestion(question, suggestions[index])
		answers = append(answers, input)
	}

	return Book{
		ID:        0,
		NAME:      answers[0],
		AUTHOR:    answers[1],
		PAGES:     answers[2],
		READCOUNT: answers[3],
		GENRE:     answers[4],
		LANGUAGE:  answers[5],
		OPINION:   answers[6],
		DATE:      answers[7],
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
