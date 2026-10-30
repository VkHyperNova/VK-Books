package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These tests reuse setupBooks, sampleBooks and loadSaved from books_test.go.
// setupBooks points LocalFile and BackupFile at a temp dir, stubs mountBackup
// to succeed and clears notices; tests that need a failing mount reassign
// mountBackup afterwards (it is restored when the test ends).

func assertSameBooks(t *testing.T, got, want *Books) {
	t.Helper()

	if len(got.Books) != len(want.Books) {
		t.Fatalf("got %d books, want %d", len(got.Books), len(want.Books))
	}
	for i := range want.Books {
		if got.Books[i] != want.Books[i] {
			t.Errorf("book %d: got %+v, want %+v", i, got.Books[i], want.Books[i])
		}
	}
}

// blockDir turns the parent directory of path into a regular file, so any
// attempt to create or write path fails.
func blockDir(t *testing.T, path string) {
	t.Helper()

	parent := filepath.Dir(path)
	if err := os.MkdirAll(filepath.Dir(parent), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(parent, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
}

/* Load */

func TestLoad_MissingFileStartsEmpty(t *testing.T) {
	local, _ := setupBooks(t)

	b := &Books{}
	if err := b.Load(local); err != nil {
		t.Fatalf("Load returned error for missing file: %v", err)
	}
	if len(b.Books) != 0 {
		t.Errorf("expected no books, got %d", len(b.Books))
	}
	if len(notices) != 1 || !strings.Contains(notices[0], "No data file found") {
		t.Errorf("expected a 'no data file' notice, got %v", notices)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	local, _ := setupBooks(t)

	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(local, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := (&Books{}).Load(local)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parsing JSON") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoad_ReadError(t *testing.T) {
	local, _ := setupBooks(t)

	// A directory can't be read as a file, and isn't ErrNotExist.
	dir := filepath.Dir(local)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := (&Books{}).Load(dir)
	if err == nil {
		t.Fatal("expected error when reading a directory")
	}
	if !strings.Contains(err.Error(), "error reading file") {
		t.Errorf("unexpected error: %v", err)
	}
}

// This pins the JSON field names, so renaming a tag can't silently stop
// existing books.json files from loading.
func TestLoad_ReadsExistingFileFormat(t *testing.T) {
	local, _ := setupBooks(t)

	raw := `{
  "books": [
    {
      "id": 1,
      "book": "Dune",
      "author": "Frank Herbert",
      "pages": "612",
      "readcount": "3",
      "genre": "sci-fi",
      "language": "English",
      "opinion": "classic",
      "date": "14.01.2026"
    }
  ]
}`
	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(local, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}

	b := &Books{}
	if err := b.Load(local); err != nil {
		t.Fatalf("Load: %v", err)
	}

	want := Book{Id: 1, Name: "Dune", Author: "Frank Herbert", Pages: "612", ReadCount: "3",
		Genre: "sci-fi", Language: "English", Opinion: "classic", Date: "14.01.2026"}
	if len(b.Books) != 1 || b.Books[0] != want {
		t.Errorf("books = %+v, want [%+v]", b.Books, want)
	}
}

/* Save */

func TestSave_RoundTrip(t *testing.T) {
	local, _ := setupBooks(t)

	want := sampleBooks()
	if err := want.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	assertSameBooks(t, loadSaved(t, local), want)
}

func TestSave_UsesExpectedJSONFieldNames(t *testing.T) {
	local, _ := setupBooks(t)

	if err := sampleBooks().Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := os.ReadFile(local)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"books"`, `"id"`, `"book"`, `"author"`, `"pages"`, `"readcount"`, `"genre"`, `"language"`, `"opinion"`, `"date"`} {
		if !strings.Contains(string(data), key) {
			t.Errorf("saved JSON should contain %s", key)
		}
	}
}

func TestSave_LocalFileIsPrivate(t *testing.T) {
	local, _ := setupBooks(t)

	if err := sampleBooks().Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(local)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("perm = %v, want 0600", info.Mode().Perm())
	}
}

func TestSave_WritesBackupWhenMounted(t *testing.T) {
	local, _ := setupBooks(t)

	if err := sampleBooks().Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	localData, err := os.ReadFile(local)
	if err != nil {
		t.Fatal(err)
	}
	backupData, err := os.ReadFile(BackupFile)
	if err != nil {
		t.Fatalf("backup not written: %v", err)
	}
	if string(localData) != string(backupData) {
		t.Error("backup contents differ from local file")
	}
	if len(notices) != 0 {
		t.Errorf("unexpected notices: %v", notices)
	}
}

func TestSave_MountFailureIsNotFatal(t *testing.T) {
	local, _ := setupBooks(t)
	mountBackup = func() error { return errors.New("drive not found") }

	if err := sampleBooks().Save(); err != nil {
		t.Fatalf("Save should succeed when backup fails, got: %v", err)
	}

	if _, err := os.Stat(local); err != nil {
		t.Errorf("local file should exist: %v", err)
	}
	if _, err := os.Stat(BackupFile); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("backup should not exist, stat err = %v", err)
	}
	if len(notices) != 1 || !strings.Contains(notices[0], "Backup failed") {
		t.Fatalf("expected a 'Backup failed' notice, got %v", notices)
	}
	if !strings.Contains(notices[0], "drive not found") {
		t.Errorf("notice should include the mount error, got %q", notices[0])
	}
}

func TestSave_BackupWriteFailureIsNotFatal(t *testing.T) {
	local, _ := setupBooks(t)
	blockDir(t, BackupFile) // mount "works", but the backup path can't be written

	if err := sampleBooks().Save(); err != nil {
		t.Fatalf("Save should succeed when the backup write fails, got: %v", err)
	}

	if _, err := os.Stat(local); err != nil {
		t.Errorf("local file should exist: %v", err)
	}
	if len(notices) != 1 || !strings.Contains(notices[0], "Backup failed") {
		t.Errorf("expected a 'Backup failed' notice, got %v", notices)
	}
}

func TestSave_LocalFailureIsFatal(t *testing.T) {
	local, _ := setupBooks(t)
	blockDir(t, local)

	err := sampleBooks().Save()
	if err == nil {
		t.Fatal("expected error when local save fails")
	}
	if !strings.Contains(err.Error(), "local save") {
		t.Errorf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(BackupFile); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("no backup should be made when the local save fails, stat err = %v", statErr)
	}
}

func TestSave_EmptyListIsValid(t *testing.T) {
	local, _ := setupBooks(t)

	if err := (&Books{}).Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if got := loadSaved(t, local); len(got.Books) != 0 {
		t.Errorf("expected no books, got %d", len(got.Books))
	}
}

func TestSave_OverwritesPreviousData(t *testing.T) {
	local, _ := setupBooks(t)

	b := sampleBooks()
	if err := b.Save(); err != nil {
		t.Fatal(err)
	}

	b.Books = b.Books[:1]
	if err := b.Save(); err != nil {
		t.Fatal(err)
	}

	if got := loadSaved(t, local); len(got.Books) != 1 {
		t.Errorf("got %d books after second save, want 1", len(got.Books))
	}
}

/* backup */

func TestBackup_WritesFileWhenMounted(t *testing.T) {
	setupBooks(t)

	if err := backup([]byte("payload")); err != nil {
		t.Fatalf("backup: %v", err)
	}

	got, err := os.ReadFile(BackupFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "payload" {
		t.Errorf("content = %q, want %q", got, "payload")
	}
}

func TestBackup_WrapsMountError(t *testing.T) {
	setupBooks(t)
	mountErr := errors.New("drive not found")
	mountBackup = func() error { return mountErr }

	err := backup([]byte("payload"))
	if !errors.Is(err, mountErr) {
		t.Errorf("expected error wrapping the mount error, got %v", err)
	}
	if _, statErr := os.Stat(BackupFile); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("backup file should not exist, stat err = %v", statErr)
	}
}

/* writeFileAtomic */

func TestWriteFileAtomic_CreatesDirsAndSetsPerm(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c.json")

	if err := writeFileAtomic(path, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Errorf("content = %q, want %q", got, "hello")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("perm = %v, want 0600", info.Mode().Perm())
	}
}

func TestWriteFileAtomic_OverwritesAndLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")

	if err := writeFileAtomic(path, []byte("first"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(path, []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second" {
		t.Errorf("content = %q, want %q", got, "second")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "data.json" {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected only data.json in dir, found %v", names)
	}
}

func TestWriteFileAtomic_FailureCleansUpTempFile(t *testing.T) {
	dir := t.TempDir()

	// Renaming a file onto a non-empty directory fails.
	target := filepath.Join(dir, "target")
	if err := os.MkdirAll(filepath.Join(target, "child"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := writeFileAtomic(target, []byte("data"), 0o600); err == nil {
		t.Fatal("expected rename onto a directory to fail")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp-") {
			t.Errorf("temp file %q was left behind", e.Name())
		}
	}
}
