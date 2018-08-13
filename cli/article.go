package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/ericaro/frontmatter"
)

//ImportArticle represents and article that can be imported
type ImportArticle struct {
	Title      string   `yaml:"title"`
	URL        string   `yaml:"url"`
	Banner     string   `yaml:"banner"`
	Date       string   `yaml:"date"`
	Author     string   `yaml:"author"`
	Categories []string `yaml:"categories"`
	Tags       []string `yaml:"tags"`
	Layout     string   `yaml:"layout"`
	Content    string   `fm:"content" yaml:"-"`
}

//LoadImportArticle load article from yaml file
func LoadImportArticle(filename string) (*ImportArticle, error) {
	importFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	importArticle := new(ImportArticle)
	err = frontmatter.Unmarshal(importFile, importArticle)

	return importArticle, err
}

//CreateOrUpdateArticle save article from input filename
func CreateOrUpdateArticle(filename string) {
	inputLocation, isFolder := getInputLocation(filename)
	if isFolder {
		subDirToSkip := []string{".git", ".DS_Store"}
		err := iterateFolder(inputLocation, "md", subDirToSkip, saveArticle)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	saveArticle(inputLocation)
}

func getInputLocation(inputLocation string) (string, bool) {
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

func saveArticle(filename string) {
	importArticle, err := LoadImportArticle(filename)
	if err != nil {
		log.Infof("Unable to save file: %v\n", filename)
		log.Error(err.Error())
		return
	}

	fmt.Printf("saving article: %v\n", importArticle.Title)
}

func iterateFolder(fileFolder, extensionToFind string, directoriesToSkip []string, action func(filename string)) error {
	err := filepath.Walk(fileFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && contains(directoriesToSkip, info.Name()) {
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

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
