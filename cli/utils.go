package cli

import (
	"bufio"
	"fmt"
	"os"
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
