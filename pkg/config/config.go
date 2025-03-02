package config

import "time"

const FileName = "books.json"
const FolderName = "BOOKS"

var Date = time.Now().Format("02.01.2006")

var LocalPath = "./" + FolderName + "/" + FileName
var BackupPath = "/media/veikko/VK DATA/DATABASES/" + FolderName + "/" + FileName
var BackupPathWithDate = "/media/veikko/VK DATA/DATABASES/" + FolderName + "/" + "(" + Date + ") " + FileName
