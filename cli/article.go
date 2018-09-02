package cli

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"net/http"

	"encoding/json"

	"bytes"

	"time"

	"io"

	"github.com/ericaro/frontmatter"
	"github.com/evcraddock/goarticles"
	"github.com/evcraddock/goarticles/configs"
	"github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2/bson"
)

//ImportArticle represents and article that can be imported
type ImportArticle struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	URL         string   `yaml:"url"`
	Banner      string   `yaml:"banner"`
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
	AuthURL     string
	Auth        AuthRequestBody
	AccessToken string
	Cache       cache.Cache
}

//AuthRequestBody request headers for getting authorization token
type AuthRequestBody struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
}

//AuthResponse response from the authorization service
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

//NewImportArticleService create new article service
func NewImportArticleService(config configs.ClientConfiguration) *ImportArticleService {
	return &ImportArticleService{
		URL:     config.URL,
		AuthURL: config.Auth.URL,
		Auth: AuthRequestBody{
			GrantType:    config.Auth.GrantType,
			ClientID:     config.Auth.ClientID,
			ClientSecret: config.Auth.ClientSecret,
			Audience:     config.Auth.Audience,
		},
		Cache: *cache.New(5*time.Minute, 10*time.Minute),
	}
}

//CreateOrUpdateArticle save article from input filename
func (s *ImportArticleService) CreateOrUpdateArticle(filename string) {
	inputLocation, isFolder := s.getInputLocation(filename)
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

func (s *ImportArticleService) getInputLocation(inputLocation string) (string, bool) {
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
		if err := s.updateArticle(*importArticle); err == nil {
			return
		}
	}

	//TODO: handle knowing if article doesn't exist better than this
	if err != nil && err.Error() == "404" {
		importArticle.ID = ""
	}

	if _, err := s.createArticle(*importArticle); err != nil {
		log.Error(err.Error())
	}
}

func (s *ImportArticleService) loadArticle(id string) (*ImportArticle, error) {
	url := s.URL + "/api/articles/" + id
	authToken := s.getAccessToken()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
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

	//json.NewDecoder(res.Body).Decode(article)
	//if err := json.NewDecoder(res.Body).Decode(article); err != nil {
	//	return nil, err
	//}

	return s.copyFrom(article)
}

func (s *ImportArticleService) updateArticle(importArticle ImportArticle) error {
	url := s.URL + "/api/articles/" + importArticle.ID
	authToken := s.getAccessToken()

	article, err := s.copyTo(&importArticle)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(article)
	client := &http.Client{}

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+authToken)
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
	authToken := s.getAccessToken()

	article, err := s.copyTo(&importArticle)
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(article)

	client := &http.Client{}

	fmt.Printf("article: %v \n", string(b))

	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 201 {
		fmt.Printf("successfully added article: %v \n", importArticle.Title)
		return article, nil
	}

	return nil, fmt.Errorf("Error updating article with status: %v", res.Status)
}

func (s *ImportArticleService) getAccessToken() string {

	if token, found := s.getTokenFromCache(); found {
		return token
	}

	b, _ := json.Marshal(s.Auth)
	res, err := http.Post(s.AuthURL, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Error("error getting token: %v", err.Error())
	}

	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	authRes := &AuthResponse{}
	if err := json.Unmarshal(resBody, authRes); err != nil {
		log.Error("error decoding body: %v", err.Error())
	}

	s.saveTokenToCache(authRes.AccessToken)

	return authRes.AccessToken
}

func (s *ImportArticleService) getTokenFromCache() (string, bool) {

	token := ""

	x, found := s.Cache.Get("token")
	if found {
		token = x.(string)
	}

	return token, found
}

func (s *ImportArticleService) saveTokenToCache(token string) {
	s.Cache.Set("token", token, cache.DefaultExpiration)
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
