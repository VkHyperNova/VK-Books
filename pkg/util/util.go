package util

import (
	"fmt"
	"os"
	"vk-books/pkg/config"

	"github.com/peterh/liner"
)

func CreateNecessaryFiles() error {

	// Make Folder
	if err := os.MkdirAll(config.FolderName, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", config.FolderName, err)
	}

	// Make File
	if _, err := os.Stat(config.LocalPath); os.IsNotExist(err) {
		if err := os.WriteFile(config.LocalPath, []byte(`{"books": []}`), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", config.FileName, err)
		}
	}

	return nil
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
