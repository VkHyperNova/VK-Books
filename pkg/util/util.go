package util

import (
	"log"
	"os"
	"vk-books/pkg/config"
)

func CreateDirectoryIfNotExists(folderName string) {
	if err := os.MkdirAll(folderName, 0700); err != nil {
		log.Println("Error creating folder:", err)
	} else {
		err = os.WriteFile(config.LocalPath, []byte(`{"books": []}`), 0644)
		if err != nil {
			panic(err)
		}
	}
}
