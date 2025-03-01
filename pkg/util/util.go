package util

import (
	"fmt"
	"os"
	"vk-books/pkg/config"
)

func CreateNecessaryFiles() error {

	// Make Folder
	if err := os.MkdirAll(config.Folder, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", config.Folder, err)
	}

	// Make File
	if _, err := os.Stat(config.LocalPath); os.IsNotExist(err) {
		if err := os.WriteFile(config.LocalPath, []byte(`{"books": []}`), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", config.File, err)
		}
	}

	return nil
}
