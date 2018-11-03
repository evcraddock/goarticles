package cli

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"

	"github.com/evcraddock/goarticles/services"
	log "github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ericaro/frontmatter"
	"github.com/evcraddock/goarticles"
	"github.com/evcraddock/goarticles/configs"
	"gopkg.in/mgo.v2/bson"
)

//ImportArticle represents and article that can be imported
type ImportArticle struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	URL         string   `yaml:"url"`
	Banner      string   `yaml:"banner"`
	Images      []string `yaml:"images"`
	PublishDate string   `yaml:"publishDate"`
	Author      string   `yaml:"author"`
	Categories  []string `yaml:"categories"`
	Tags        []string `yaml:"tags"`
	Layout      string   `yaml:"layout"`
	Content     string   `fm:"content" yaml:"-"`
}

//ImportArticleService service used to handle interactions with the API
type ImportArticleService struct {
	URL         string
	AccessToken string
}

//NewImportArticleService create new article service
func NewImportArticleService(config configs.ClientConfiguration) *ImportArticleService {

	accessTokenService := services.NewAccessTokenService(config.Auth.URL,
		services.AuthRequestBody{
			GrantType:    config.Auth.GrantType,
			ClientID:     config.Auth.ClientID,
			ClientSecret: config.Auth.ClientSecret,
			Audience:     config.Auth.Audience,
		},
	)

	return &ImportArticleService{
		URL:         config.URL,
		AccessToken: accessTokenService.GetAccessToken(),
	}
}

//CreateOrUpdateArticle save article from input filename
func (s *ImportArticleService) CreateOrUpdateArticle(filename string) {
	inputLocation, isFolder := GetInputLocation(filename)
	if isFolder {
		subDirToSkip := []string{".git", ".DS_Store"}
		err := IterateFolder(inputLocation, "md", subDirToSkip, s.saveArticle)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	s.saveArticle(inputLocation)
}

func (s *ImportArticleService) loadImportArticle(filename string) (*ImportArticle, error) {
	importFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	importArticle := new(ImportArticle)
	err = frontmatter.Unmarshal(importFile, importArticle)

	return importArticle, err
}

func (s *ImportArticleService) saveArticle(filename string) {
	importArticle, err := s.loadImportArticle(filename)
	if err != nil {
		log.Debugf("Unable to save file: %v\n", filename)
		log.Error(err.Error())
		return
	}

	if importArticle.ID != "" {
		if err := s.updateArticle(*importArticle); err != nil {
			if err.Error() == "404" {
				importArticle.ID = ""
			}
		}
	}

	savedArticleID := importArticle.ID

	if savedArticleID == "" {
		newArticle, err := s.createArticle(*importArticle)
		if err != nil {
			log.Error(err.Error())
			return
		}

		savedArticleID = newArticle.ID.Hex()
	}

	if len(importArticle.Images) > 0 {
		fileDir := filepath.Dir(filename)
		if err := s.saveImages(savedArticleID, fileDir, importArticle.Images); err != nil {
			log.Error(err.Error())
		}
	}

	importArticle.ID = savedArticleID

	if err := s.saveMarkdownFile(importArticle, filename); err != nil {
		log.Error(err.Error())
	}

}

func (s *ImportArticleService) saveMarkdownFile(importArticle *ImportArticle, filename string) error {

	data, err := frontmatter.Marshal(importArticle)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *ImportArticleService) saveImages(id string, directory string, images []string) error {
	for _, filename := range images {
		if err := s.saveImage(id, directory, filename); err != nil {
			log.Error(err.Error())
		}
	}

	return nil
}

func (s *ImportArticleService) saveImage(id string, directory string, filename string) error {
	url := s.URL + "/api/articles/" + id + "/images"

	client := &http.Client{}
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, filename))

	file, err := os.Open(directory + "/" + filename)
	if err != nil {
		return err
	}

	defer file.Close()

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	fileWriter, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	writer.Close()

	req, _ := http.NewRequest("POST", url, buffer)
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 202 {
		fmt.Printf("successfully added article image: %v \n", filename)
		return nil
	}

	return fmt.Errorf("failed to save image with error: %v", res.Status)
}

func (s *ImportArticleService) loadArticle(id string) (*ImportArticle, error) {
	url := s.URL + "/api/articles/" + id

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error: %v", res.Status)
	}

	article := &goarticles.Article{}

	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err := json.Unmarshal(body, &article); err != nil {
		return nil, err
	}

	return s.copyFrom(article)
}

func (s *ImportArticleService) updateArticle(importArticle ImportArticle) error {
	url := s.URL + "/api/articles/" + importArticle.ID

	article, err := s.copyTo(&importArticle)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(article)
	client := &http.Client{}

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Printf("successfully updated article: %v \n", importArticle.Title)
		return nil
	}

	return fmt.Errorf("%v", res.StatusCode)
}

func (s *ImportArticleService) createArticle(importArticle ImportArticle) (*goarticles.Article, error) {
	url := s.URL + "/api/articles"

	article, err := s.copyTo(&importArticle)
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(article)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 201 {

		body, _ := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
		if err := json.Unmarshal(body, &article); err != nil {
			return nil, err
		}

		fmt.Printf("successfully added article: %v \n", article.Title)
		return article, nil
	}

	return nil, fmt.Errorf("Error updating article with status: %v", res.Status)
}

func (s *ImportArticleService) copyFrom(article *goarticles.Article) (*ImportArticle, error) {
	importArticle := &ImportArticle{
		ID:          article.ID.Hex(),
		Title:       article.Title,
		URL:         article.URL,
		Author:      article.Author,
		Banner:      article.Banner,
		Categories:  article.Categories,
		Content:     article.Content,
		PublishDate: article.PublishDate.Format("2006-01-02"),
		Tags:        article.Tags,
	}

	return importArticle, nil
}

func (s *ImportArticleService) copyTo(importArticle *ImportArticle) (*goarticles.Article, error) {
	article := &goarticles.Article{
		Title:      importArticle.Title,
		URL:        importArticle.URL,
		Author:     importArticle.Author,
		Banner:     importArticle.Banner,
		Categories: importArticle.Categories,
		Content:    importArticle.Content,
		Tags:       importArticle.Tags,
	}

	if importArticle.ID != "" {
		article.ID = bson.ObjectIdHex(importArticle.ID)
	}

	importPublishDate, err := time.Parse("01/02/2006", importArticle.PublishDate)
	if err != nil {
		return nil, err
	}

	article.PublishDate = importPublishDate

	return article, nil
}
