package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"vk-books/pkg/color"
	"vk-books/pkg/config"

	"github.com/peterh/liner"
)

func AutoDetectLanguage(name string) string {

	for _, char := range name {
		if unicode.In(char, unicode.Cyrillic) {
			return "Russian"
		}
	}

	return "English"
}

func ensureFile(path string, content string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creating directory for %s: %w", path, err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", path, err)
		}
	}

	return nil
}

func CreateFilesAndFolders() error {
	
	if err := ensureFile(config.LocalFile, config.DefaultContent); err != nil {
		return err
	}

	if !HardDriveMountCheck() {
		input := Prompt("Do you want to continue? (y/n) ")
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
	} else {
		if err := ensureFile(config.BackupFile, config.DefaultContent); err != nil {
			return err
		}
	}

	return nil
}

func HardDriveMountCheck() bool {
	if runtime.GOOS != "linux" {
		fmt.Println("This program only works on Linux.")
		return false
	}

	mountPoint := "/media/veikko/VK\\040DATA" // match /proc/mounts format

	file, err := os.Open("/proc/mounts")
	if err != nil {
		fmt.Println("Cannot open /proc/mounts:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[1] == mountPoint {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning /proc/mounts:", err)
		return false
	}

	fmt.Println(color.Red + "\nVK DATA is NOT mounted" + color.Reset)
	return false
}

func Prompt(Question string) string {

	fmt.Print(color.Cyan, Question, color.Reset)

	var input string

	fmt.Scanln(&input)

	return input
}

func PromptWithSuggestion(name string, suggestion string) (string, error) {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", suggestion, -1)
	if err != nil {
		return input, err
	}

	return input, nil
}

func ClearScreen() {

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}

func Contains(arr []string, item string) bool {
	for _, str := range arr {
		if str == item {
			return true
		}
	}
	return false
}

func Confirm() bool {

	input := Prompt("(y/n): ")

	if input == "n" || input == "no" || input == "q" {
		fmt.Println(color.Red, "Aborted!", color.Reset)
		return false
	}
	return true
}