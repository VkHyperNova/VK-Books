package util

import (
	"fmt"
	"os"
)

func CreateDirectory(folderName string) {
	if err := os.MkdirAll(folderName, 0700); err != nil {
		fmt.Println("Error creating folder:", err)
	} else {
		fmt.Println(folderName, "folder is ready")
	}
}
