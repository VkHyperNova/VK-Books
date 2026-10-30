package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

/* Test helpers */

// cancelAnswer makes the stubbed prompt behave like Ctrl+C.
const cancelAnswer = "\x00cancel"

type promptCall struct{ label, suggestion string }

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string { return ansiRe.ReplaceAllString(s, "") }

// setupBooks points the storage globals at a temp dir, stubs the drive mount
// and replaces prompt with a stub that hands out the given answers in order.
// Everything is restored when the test finishes.
func setupBooks(t *testing.T, answers ...string) (local string, calls *[]promptCall) {
	t.Helper()

	dir := t.TempDir()
	local = filepath.Join(dir, "local", "books.json")
	backup := filepath.Join(dir, "backup", "books.json")

	oldLocal, oldBackup, oldMount, oldNotices, oldPrompt := LocalFile, BackupFile, mountBackup, notices, prompt

	calls = &[]promptCall{}
	remaining := answers

	LocalFile, BackupFile = local, backup
	mountBackup = func() error { return nil }
	notices = nil
	prompt = func(label, suggestion string) (string, error) {
		*calls = append(*calls, promptCall{label, suggestion})
		if len(remaining) == 0 {
			return "", fmt.Errorf("unexpected prompt %q", label)
		}
		a := remaining[0]
		remaining = remaining[1:]
		if a == cancelAnswer {
			return "", errCancelled
		}
		return a, nil
	}

	t.Cleanup(func() {
		LocalFile, BackupFile, mountBackup, notices, prompt = oldLocal, oldBackup, oldMount, oldNotices, oldPrompt
	})
	return local, calls
}

// setArgs replaces os.Args for the duration of the test.
func setArgs(t *testing.T, args ...string) {
	t.Helper()

	old := os.Args
	os.Args = args
	t.Cleanup(func() { os.Args = old })
}

// captureStdout returns everything f printed to stdout.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = old })

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		done <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = old
	return <-done
}

func findCall(t *testing.T, calls []promptCall, label string) promptCall {
	t.Helper()
	for _, c := range calls {
		if c.label == label {
			return c
		}
	}
	t.Fatalf("no prompt with label %q", label)
	return promptCall{}
}

func sampleBooks() *Books {
	return &Books{Books: []Book{
		{Id: 1, Name: "Dune", Author: "Frank Herbert", Pages: "612", ReadCount: "3", Genre: "sci-fi", Language: "English", Opinion: "classic", Date: "14.01.2026"},
		{Id: 2, Name: "Neuromancer", Author: "William Gibson", Pages: "271", ReadCount: "1", Genre: "sci-fi", Language: "English", Opinion: "dense", Date: "30.04.2026"},
		{Id: 3, Name: "Dune Messiah", Author: "Frank Herbert", Pages: "256", ReadCount: "1", Genre: "sci-fi", Language: "English", Opinion: "", Date: "02.05.2026"},
	}}
}

func loadSaved(t *testing.T, path string) *Books {
	t.Helper()

	b := &Books{}
	if err := b.Load(path); err != nil {
		t.Fatalf("load saved file: %v", err)
	}
	return b
}

func assertNotSaved(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("nothing should have been saved, stat err = %v", err)
	}
}

func ids(b *Books) []int {
	out := make([]int, len(b.Books))
	for i, bk := range b.Books {
		out[i] = bk.Id
	}
	return out
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

/* parseID, indexOf, nextID, detectLanguage */

func TestParseID(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"3", 3, false},
		{" 7 ", 7, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"1.5", 0, true},
		{"", 0, true},
	}
	for _, tt := range tests {
		got, err := parseID(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseID(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("parseID(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestIndexOf(t *testing.T) {
	b := sampleBooks()

	if i, err := b.indexOf(2); err != nil || i != 1 {
		t.Errorf("indexOf(2) = %d, %v; want 1, nil", i, err)
	}

	if _, err := b.indexOf(99); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("indexOf(99) error = %v, want 'not found'", err)
	}

	for _, id := range []int{0, -1} {
		if _, err := b.indexOf(id); err == nil || !strings.Contains(err.Error(), "invalid ID") {
			t.Errorf("indexOf(%d) error = %v, want 'invalid ID'", id, err)
		}
	}
}

func TestNextID(t *testing.T) {
	tests := []struct {
		name string
		ids  []int
		want int
	}{
		{"empty", nil, 1},
		{"sequential", []int{1, 2, 3}, 4},
		{"gaps", []int{1, 2, 5}, 6},
		{"unsorted", []int{5, 1, 2}, 6},
	}
	for _, tt := range tests {
		b := &Books{}
		for _, id := range tt.ids {
			b.Books = append(b.Books, Book{Id: id})
		}
		if got := b.nextID(); got != tt.want {
			t.Errorf("%s: nextID() = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Dune", "English"},
		{"Война и мир", "Russian"},
		{"Dune Война", "Russian"}, // any Cyrillic letter wins
		{"", "English"},
	}
	for _, tt := range tests {
		if got := detectLanguage(tt.in); got != tt.want {
			t.Errorf("detectLanguage(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

/* Validators */

func TestValidators(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(string) error
		in      string
		wantErr bool
	}{
		{"required ok", required, "Dune", false},
		{"required empty", required, "", true},

		{"positiveInt 1", positiveInt, "1", false},
		{"positiveInt 612", positiveInt, "612", false},
		{"positiveInt 0", positiveInt, "0", true},
		{"positiveInt negative", positiveInt, "-3", true},
		{"positiveInt text", positiveInt, "abc", true},
		{"positiveInt decimal", positiveInt, "1.5", true},
		{"positiveInt empty", positiveInt, "", true},

		{"validDate ok", validDate, "14.01.2026", false},
		{"validDate iso", validDate, "2026-01-14", true},
		{"validDate impossible", validDate, "31.02.2026", true},
		{"validDate empty", validDate, "", true},
		{"validDate text", validDate, "yesterday", true},
	}
	for _, tt := range tests {
		if err := tt.fn(tt.in); (err != nil) != tt.wantErr {
			t.Errorf("%s: error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

/* Add */

func TestAdd_AddsFirstBookWithIDAndSaves(t *testing.T) {
	local, _ := setupBooks(t, "Dune", "Frank Herbert", "612", "1", "sci-fi", "English", "great", "14.01.2026")

	b := &Books{}
	if err := b.Add(); err != nil {
		t.Fatalf("Add: %v", err)
	}

	want := Book{Id: 1, Name: "Dune", Author: "Frank Herbert", Pages: "612", ReadCount: "1",
		Genre: "sci-fi", Language: "English", Opinion: "great", Date: "14.01.2026"}

	if len(b.Books) != 1 || b.Books[0] != want {
		t.Fatalf("books = %+v, want [%+v]", b.Books, want)
	}

	saved := loadSaved(t, local)
	if len(saved.Books) != 1 || saved.Books[0] != want {
		t.Errorf("saved = %+v, want [%+v]", saved.Books, want)
	}
}

func TestAdd_AssignsNextID(t *testing.T) {
	setupBooks(t, "Hyperion", "Dan Simmons", "482", "1", "sci-fi", "English", "", "01.06.2026")

	b := sampleBooks()
	if err := b.Add(); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if len(b.Books) != 4 {
		t.Fatalf("got %d books, want 4", len(b.Books))
	}
	if last := b.Books[3]; last.Id != 4 || last.Name != "Hyperion" {
		t.Errorf("last book = %+v, want id 4 'Hyperion'", last)
	}
}

func TestAdd_SuggestsDefaults(t *testing.T) {
	_, calls := setupBooks(t, "Dune", "Frank Herbert", "612", "1", "sci-fi", "English", "", "14.01.2026")

	if err := (&Books{}).Add(); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if got := findCall(t, *calls, "Pages Count:").suggestion; got != "1" {
		t.Errorf("pages suggestion = %q, want %q", got, "1")
	}
	if got := findCall(t, *calls, "Read Count:").suggestion; got != "1" {
		t.Errorf("read count suggestion = %q, want %q", got, "1")
	}
	if got, want := findCall(t, *calls, "Date:").suggestion, time.Now().Format(dateLayout); got != want {
		t.Errorf("date suggestion = %q, want today (%q)", got, want)
	}
}

func TestAdd_LanguageGuessUsesTheTypedName(t *testing.T) {
	_, calls := setupBooks(t, "Война и мир", "Лев Толстой", "1225", "1", "classic", "Russian", "", "14.01.2026")

	if err := (&Books{}).Add(); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if got := findCall(t, *calls, "Language:").suggestion; got != "Russian" {
		t.Errorf("language suggestion = %q, want %q", got, "Russian")
	}
}

func TestAdd_InvalidInputIsAskedAgainWithTypedValue(t *testing.T) {
	local, calls := setupBooks(t,
		"Dune", "Frank Herbert",
		"abc", "300", // invalid pages, then a valid one
		"1", "sci-fi", "English", "", "14.01.2026")

	b := &Books{}
	var err error
	captureStdout(t, func() { err = b.Add() })
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	// name, author, pages (abc), pages again
	if len(*calls) < 4 {
		t.Fatalf("expected at least 4 prompts, got %d", len(*calls))
	}
	retry := (*calls)[3]
	if retry.label != "Pages Count:" || retry.suggestion != "abc" {
		t.Errorf("retry prompt = %+v, want pages prompt suggesting the typed %q", retry, "abc")
	}
	if b.Books[0].Pages != "300" {
		t.Errorf("pages = %q, want %q", b.Books[0].Pages, "300")
	}
	if saved := loadSaved(t, local); saved.Books[0].Pages != "300" {
		t.Errorf("saved pages = %q, want %q", saved.Books[0].Pages, "300")
	}
}

func TestAdd_RequiresNameAndAuthor(t *testing.T) {
	_, calls := setupBooks(t,
		"", "Dune", // empty name, then a valid one
		"", "Frank Herbert", // empty author, then a valid one
		"612", "1", "sci-fi", "English", "", "14.01.2026")

	b := &Books{}
	var err error
	captureStdout(t, func() { err = b.Add() })
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if b.Books[0].Name != "Dune" || b.Books[0].Author != "Frank Herbert" {
		t.Errorf("book = %+v", b.Books[0])
	}
	if len(*calls) != 10 {
		t.Errorf("got %d prompts, want 10 (8 fields + 2 retries)", len(*calls))
	}
}

func TestAdd_CancelSavesNothing(t *testing.T) {
	local, _ := setupBooks(t, "Dune", cancelAnswer)

	b := &Books{}
	err := b.Add()

	if !errors.Is(err, errCancelled) {
		t.Fatalf("error = %v, want errCancelled", err)
	}
	if len(b.Books) != 0 {
		t.Errorf("book was added despite cancel: %+v", b.Books)
	}
	assertNotSaved(t, local)
}

/* Update */

func TestUpdate_ChangesOnlyTheChosenBookAndSaves(t *testing.T) {
	local, calls := setupBooks(t, "New Name", "Frank", "300", "2", "fantasy", "English", "meh", "01.02.2026")
	setArgs(t, "prog", "update", "2")

	b := sampleBooks()
	orig := sampleBooks()
	if err := b.Update(); err != nil {
		t.Fatalf("Update: %v", err)
	}

	want := Book{Id: 2, Name: "New Name", Author: "Frank", Pages: "300", ReadCount: "2",
		Genre: "fantasy", Language: "English", Opinion: "meh", Date: "01.02.2026"}

	if b.Books[1] != want {
		t.Errorf("updated book = %+v, want %+v", b.Books[1], want)
	}
	if b.Books[0] != orig.Books[0] || b.Books[2] != orig.Books[2] {
		t.Error("other books were changed")
	}

	// The prompts start from the existing values
	if got := findCall(t, *calls, "Book Name:").suggestion; got != "Neuromancer" {
		t.Errorf("name suggestion = %q, want %q", got, "Neuromancer")
	}
	if got := findCall(t, *calls, "Pages Count:").suggestion; got != "271" {
		t.Errorf("pages suggestion = %q, want %q", got, "271")
	}

	saved := loadSaved(t, local)
	if len(saved.Books) != 3 || saved.Books[1] != want {
		t.Errorf("saved = %+v", saved.Books)
	}
}

func TestUpdate_ArgumentErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{"missing id", []string{"prog", "update"}, "usage: update"},
		{"not a number", []string{"prog", "update", "abc"}, "invalid ID"},
		{"unknown id", []string{"prog", "update", "99"}, "not found"},
	}
	for _, tt := range tests {
		local, calls := setupBooks(t)
		setArgs(t, tt.args...)

		b := sampleBooks()
		err := b.Update()
		if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
			t.Errorf("%s: error = %v, want it to contain %q", tt.name, err, tt.wantErr)
		}
		if len(*calls) != 0 {
			t.Errorf("%s: prompted despite bad arguments", tt.name)
		}
		assertNotSaved(t, local)
	}
}

func TestUpdate_CancelLeavesBookUnchanged(t *testing.T) {
	local, _ := setupBooks(t, "Changed Name", cancelAnswer)
	setArgs(t, "prog", "update", "2")

	b := sampleBooks()
	err := b.Update()

	if !errors.Is(err, errCancelled) {
		t.Fatalf("error = %v, want errCancelled", err)
	}
	if b.Books[1].Name != "Neuromancer" {
		t.Errorf("name = %q, want it unchanged", b.Books[1].Name)
	}
	assertNotSaved(t, local)
}

/* Delete */

func TestDelete_ConfirmedRemovesBookAndSaves(t *testing.T) {
	for _, answer := range []string{"y", "yes", "YES", "Y"} {
		local, _ := setupBooks(t, answer)
		setArgs(t, "prog", "delete", "2")

		b := sampleBooks()
		var err error
		captureStdout(t, func() { err = b.Delete() })

		if err != nil {
			t.Fatalf("answer %q: Delete: %v", answer, err)
		}
		if got := ids(b); !sameInts(got, []int{1, 3}) {
			t.Errorf("answer %q: ids = %v, want [1 3]", answer, got)
		}
		if got := ids(loadSaved(t, local)); !sameInts(got, []int{1, 3}) {
			t.Errorf("answer %q: saved ids = %v, want [1 3]", answer, got)
		}
	}
}

func TestDelete_KeepsOrderOfRemainingBooks(t *testing.T) {
	setupBooks(t, "y")
	setArgs(t, "prog", "delete", "1")

	b := sampleBooks()
	captureStdout(t, func() { b.Delete() })

	if got := ids(b); !sameInts(got, []int{2, 3}) {
		t.Errorf("ids = %v, want [2 3]", got)
	}
}

func TestDelete_DeclinedKeepsBookAndIsNotAnError(t *testing.T) {
	for _, answer := range []string{"n", "", "no", "maybe"} {
		local, _ := setupBooks(t, answer)
		setArgs(t, "prog", "delete", "2")

		b := sampleBooks()
		var err error
		out := captureStdout(t, func() { err = b.Delete() })

		if err != nil {
			t.Errorf("answer %q: declining should not be an error, got %v", answer, err)
		}
		if len(b.Books) != 3 {
			t.Errorf("answer %q: got %d books, want 3", answer, len(b.Books))
		}
		if !strings.Contains(stripANSI(out), "Not deleted.") {
			t.Errorf("answer %q: output should say 'Not deleted.', got %q", answer, out)
		}
		assertNotSaved(t, local)
	}
}

func TestDelete_ArgumentErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{"missing id", []string{"prog", "delete"}, "usage: delete"},
		{"not a number", []string{"prog", "delete", "abc"}, "invalid ID"},
		{"unknown id", []string{"prog", "delete", "99"}, "not found"},
	}
	for _, tt := range tests {
		local, calls := setupBooks(t)
		setArgs(t, tt.args...)

		b := sampleBooks()
		err := b.Delete()
		if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
			t.Errorf("%s: error = %v, want it to contain %q", tt.name, err, tt.wantErr)
		}
		if len(*calls) != 0 {
			t.Errorf("%s: asked for confirmation despite bad arguments", tt.name)
		}
		if len(b.Books) != 3 {
			t.Errorf("%s: books were removed", tt.name)
		}
		assertNotSaved(t, local)
	}
}

func TestDelete_CancelKeepsBook(t *testing.T) {
	local, _ := setupBooks(t, cancelAnswer)
	setArgs(t, "prog", "delete", "2")

	b := sampleBooks()
	var err error
	captureStdout(t, func() { err = b.Delete() })

	if !errors.Is(err, errCancelled) {
		t.Fatalf("error = %v, want errCancelled", err)
	}
	if len(b.Books) != 3 {
		t.Errorf("got %d books, want 3", len(b.Books))
	}
	assertNotSaved(t, local)
}

/* Search */

func TestSearch_FindsPartialCaseInsensitiveMatches(t *testing.T) {
	b := sampleBooks()

	var err error
	out := stripANSI(captureStdout(t, func() { err = b.Search("  DUN ") }))

	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	for _, want := range []string{"Dune", "Dune Messiah", "#1", "#3", "2 found"} {
		if !strings.Contains(out, want) {
			t.Errorf("output should contain %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Neuromancer") {
		t.Errorf("output should not contain Neuromancer:\n%s", out)
	}
}

func TestSearch_SingleMatch(t *testing.T) {
	b := sampleBooks()

	out := stripANSI(captureStdout(t, func() { b.Search("neuro") }))

	if !strings.Contains(out, "Neuromancer") || !strings.Contains(out, "1 found") {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestSearch_NoMatches(t *testing.T) {
	b := sampleBooks()

	var err error
	out := stripANSI(captureStdout(t, func() { err = b.Search("zzz") }))

	if err != nil {
		t.Errorf("no matches should not be an error, got %v", err)
	}
	if !strings.Contains(out, `No books found for "zzz"`) {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestSearch_EmptyQueryIsAnError(t *testing.T) {
	for _, q := range []string{"", "   "} {
		err := sampleBooks().Search(q)
		if err == nil || !strings.Contains(err.Error(), "usage: search") {
			t.Errorf("query %q: error = %v, want usage error", q, err)
		}
	}
}

func TestSearch_SearchesNameOnly(t *testing.T) {
	b := sampleBooks()

	out := stripANSI(captureStdout(t, func() { b.Search("gibson") })) // an author, not a title

	if !strings.Contains(out, "No books found") {
		t.Errorf("author names should not match, got:\n%s", out)
	}
}
