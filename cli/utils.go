package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

//GetInputLocation get folder location and return if it exists
func GetInputLocation(inputLocation string) (string, bool) {
	label := "Please enter file or folder name"

	if len(inputLocation) == 0 {
		inputLocation = InputPrompt(label, true)
	}

	ok, err := IsValidFolder(inputLocation)
	if !ok {
		if err != nil {
			fmt.Printf("Not a valid file or folder. \n")
			return InputPrompt(label, true), ok
		}
	}

	return inputLocation, ok
}

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

//IterateFolder run action for each folder recursively
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

//Contains check if a value is contained in an array
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
