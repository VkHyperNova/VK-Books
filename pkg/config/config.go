package config

import (
	"path/filepath"
	"time"
)
var Date = time.Now().Format("02.01.2006")

var DefaultContent = `{"books": []}`
var	file = "books.json"
var BaseDB = "BOOKS"
var BaseLocal = "DATABASES"
var	BaseBackup = "/media/veikko/VK DATA/"

var LocalFile = filepath.Join(BaseLocal, BaseDB, file)
var BackupFile = filepath.Join(BaseBackup, BaseLocal, BaseDB, file)
var BackupFileWithDate = filepath.Join(BaseBackup, BaseLocal, BaseDB, "books " + Date + ".json")




var Questions = []string {"Book:", "Author:", "Pages:", "Read Count:", "Genre:", "Language:", "Opinion:", "Date:"}
var AddSuggestions = []string{"", "", "1", "1", "Unknown","English", "", Date}

