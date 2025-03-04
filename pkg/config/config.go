package config

import "time"

const FileName = "books.json"
const FolderName = "BOOKS"

var Date = time.Now().Format("02.01.2006")

var Questions = []string {"Book:", "Author:", "Pages:", "Read Count:", "Genre:", "Language:", "Opinion:", "Date:"}
var AddSuggestions = []string{"", "", "1", "1", "Unknown","English", "", Date}

var LocalPath = "./" + FolderName + "/" + FileName
var BackupPath = "/media/veikko/VK DATA/DATABASES/" + FolderName + "/" + FileName
var BackupPathWithDate = "/media/veikko/VK DATA/DATABASES/" + FolderName + "/" + "(" + Date + ") " + FileName

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)
