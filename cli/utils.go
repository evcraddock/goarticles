package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

//InputPrompt ask for a input value
func InputPrompt(label string, required bool) string {
	input := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s : \n", label)
	for input.Scan() {

		inputValue := input.Text()
		if !required || len(inputValue) > 0 {
			return inputValue
		}

		fmt.Printf("%s : \n", label)
	}

	return ""
}

//IsValidFolder checks to see if folder exists
func IsValidFolder(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	fileInfo.Mode()
	return fileInfo.IsDir(), err
}

func IterateFolder(fileFolder, extensionToFind string, directoriesToSkip []string, action func(filename string)) error {
	err := filepath.Walk(fileFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && Contains(directoriesToSkip, info.Name()) {
			return filepath.SkipDir
		}

		if info.IsDir() == false {
			filename := info.Name()
			extension := filepath.Ext(filename)

			if extension == "."+extensionToFind {
				action(path)
			}
		}

		return err
	})

	return err
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
