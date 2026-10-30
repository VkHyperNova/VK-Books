package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/peterh/liner"
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

/* Main Commands */

func (b *Books) Add() error {

	newBook, err := b.promptBookInput(Book{})
	if err != nil {
		return err
	}

	// Add unique ID
	newBook.Id = b.nextID()

	// Add
	b.Books = append(b.Books, newBook)

	return b.Save()
}

func (b *Books) Update() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: update <id>")
	}

	id, err := parseID(os.Args[2])
	if err != nil {
		return err
	}

	index, err := b.indexOf(id)
	if err != nil {
		return err
	}

	updated, err := b.promptBookInput(b.Books[index])
	if err != nil {
		return err
	}

	b.Books[index] = updated

	return b.Save()
}

func (b *Books) Delete() error {

	if len(os.Args) < 3 {
		return fmt.Errorf("usage: delete <id>")
	}

	id, err := parseID(os.Args[2])
	if err != nil {
		return err
	}

	index, err := b.indexOf(id)
	if err != nil {
		return err
	}

	bk := b.Books[index]
	fmt.Printf("\nDelete %s#%d%s %s%s%s - %s?\n",
		Gray, bk.Id, Reset, Yellow, bk.Name, Reset, bk.Author)

	answer, err := prompt("(y/N):", "")
	if err != nil {
		return err
	}
	if a := strings.ToLower(answer); a != "y" && a != "yes" {
		fmt.Println(Gray + "Not deleted." + Reset)
		return nil
	}

	b.Books = slices.Delete(b.Books, index, index+1)
	return b.Save()
}

func (b *Books) Search(query string) error {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return fmt.Errorf("usage: search <name>")
	}

	var matches []Book
	for _, bk := range b.Books {
		if strings.Contains(strings.ToLower(bk.Name), query) {
			matches = append(matches, bk)
		}
	}

	fmt.Println(Cyan + Bold + "------------------------" + Reset)
	fmt.Println(Cyan + Bold + "SEARCH" + Reset)
	fmt.Println(Cyan + Bold + "------------------------" + Reset)

	if len(matches) == 0 {
		fmt.Printf("\n%sNo books found for %q%s\n", Gray, query, Reset)
		return nil
	}

	const maxTitle, maxAuthor = 32, 22

	idW, titleW, authorW := 0, 0, 0
	for _, bk := range matches {
		idW = max(idW, len(strconv.Itoa(bk.Id)))
		titleW = max(titleW, utf8.RuneCountInString(truncate(bk.Name, maxTitle)))
		authorW = max(authorW, utf8.RuneCountInString(truncate(bk.Author, maxAuthor)))
	}

	fmt.Println()
	for _, bk := range matches {
		fmt.Printf("\t%s#%-*d%s  %s%s%s  %s%s%s  %s%s%s\n",
			Gray, idW, bk.Id, Reset,
			Yellow, padRight(truncate(bk.Name, maxTitle), titleW), Reset,
			Gray, padRight(truncate(bk.Author, maxAuthor), authorW), Reset,
			Gray, bk.Date, Reset)
	}

	fmt.Printf("\n%s%d found%s\n", Bold, len(matches), Reset)
	return nil
}

/* Helpers */

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid ID: %q", s)
	}
	return id, nil
}

func (b *Books) indexOf(id int) (int, error) {
	if id <= 0 {
		return -1, fmt.Errorf("invalid ID: %d", id)
	}
	for i, bk := range b.Books {
		if bk.Id == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
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

func detectLanguage(name string) string {

	for _, char := range name {
		if unicode.In(char, unicode.Cyrillic) {
			return "Russian"
		}
	}

	return "English"
}

func clearScreen() {

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}

/* Prompt helpers */

var errCancelled = errors.New("cancelled")

var line *liner.State

// prompt is a variable, not a plain function, so tests can replace it
// with a fake that returns scripted answers instead of reading the keyboard.
var prompt = func(label, suggestion string) (string, error) {
	if line == nil {
		line = liner.NewLiner()
		line.SetCtrlCAborts(true)
	}
	input, err := line.PromptWithSuggestion("   "+label+" ", suggestion, -1)
	if errors.Is(err, liner.ErrPromptAborted) || errors.Is(err, io.EOF) {
		return "", errCancelled // Ctrl+C / Ctrl+D
	}
	return strings.TrimSpace(input), err
}

func closePrompt() {
	if line != nil {
		line.Close()
		line = nil
	}
}

/* Input */

const dateLayout = "02.01.2006"

func required(s string) error {
	if s == "" {
		return errors.New("required")
	}
	return nil
}

func positiveInt(s string) error {
	if n, err := strconv.Atoi(s); err != nil || n <= 0 {
		return errors.New("enter a whole number above 0")
	}
	return nil
}

func validDate(s string) error {
	if _, err := time.Parse(dateLayout, s); err != nil {
		return errors.New("use the format dd.mm.yyyy")
	}
	return nil
}

func (b *Books) promptBookInput(s Book) (Book, error) {
	if s.Pages == "" {
		s.Pages = "1"
	}
	if s.ReadCount == "" {
		s.ReadCount = "1"
	}
	if s.Date == "" {
		s.Date = time.Now().Format(dateLayout)
	}

	fields := []struct {
		label    string
		target   *string
		validate func(string) error
	}{
		{"Book Name:", &s.Name, required},
		{"Author:", &s.Author, required},
		{"Pages Count:", &s.Pages, positiveInt},
		{"Read Count:", &s.ReadCount, positiveInt},
		{"Genre:", &s.Genre, nil},
		{"Language:", &s.Language, nil},
		{"Opinion:", &s.Opinion, nil},
		{"Date:", &s.Date, validDate},
	}

	for _, f := range fields {
		// The language guess depends on the name, so make it once the name is known
		if f.target == &s.Language && s.Language == "" {
			s.Language = detectLanguage(s.Name)
		}

		suggestion := *f.target
		for {
			val, err := prompt(f.label, suggestion)
			if err != nil {
				return Book{}, err
			}
			if f.validate != nil {
				if verr := f.validate(val); verr != nil {
					fmt.Printf("   %s%v%s\n", Red, verr, Reset)
					suggestion = val // keep what they typed so they can fix it
					continue
				}
			}
			*f.target = val
			break
		}
	}
	return s, nil
}
