package util

import (
	"fmt"
	"os"
	"vk-books/pkg/config"
)

func CreateDirectoryIfNotExists(folderName string) {
	if _, err := os.Stat(config.LocalPath); os.IsNotExist(err) {

		// Create dir
		err := os.Mkdir(folderName, 0700)
		if err != nil {
			fmt.Println(err)
		}

		// Create file
		err = os.WriteFile(config.LocalPath, []byte(`{"books": []}`), 0644)
		if err != nil {
			panic(err)
		}
	}
}


