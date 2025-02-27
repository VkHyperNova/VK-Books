package config

import "time"

const FileName = "books.json"

var CurrentDate = time.Now().AddDate(0, -1, 0).Format("02.01.2006")

var LocalPath = "./BOOKS/" + FileName

var BackupPath = "/media/veikko/VK DATA/DATABASES/BOOKS/" + CurrentDate + "-" + FileName
