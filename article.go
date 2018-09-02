package goarticles

import (
	"encoding/json"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Article model
type Article struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	URL         string        `json:"url"`
	Content     string        `json:"content"`
	Banner      string        `json:"banner"`
	DataSource  string        `json:"dataSource"`
	PublishDate time.Time     `json:"publishDate"`
	Categories  []string      `json:"categories"`
	Tags        []string      `json:"tags"`
}

//Articles collection of articles
type Articles []Article

func (article *Article) MarshalJSON() ([]byte, error) {
	id := ""
	if article.ID.Valid() {
		id = article.ID.Hex()
	}

	date := ""
	if !article.PublishDate.IsZero() {
		date = article.PublishDate.Format("2006-01-02")
	}

	type Alias Article
	return json.Marshal(&struct {
		ID          string `json:"id,omitempty"`
		PublishDate string `json:"publishDate,omitempty"`
		*Alias
	}{
		ID:          id,
		PublishDate: date,
		Alias:       (*Alias)(article),
	})
}

func (article *Article) UnmarshalJSON(data []byte) error {
	id := ""
	if article.ID.Valid() {
		id = article.ID.Hex()
	}

	date := ""
	if !article.PublishDate.IsZero() {
		date = article.PublishDate.Format("2006-01-02")
	}

	type Alias Article
	aux := &struct {
		ID          string `json:"id,omitempty"`
		PublishDate string `json:"publishDate,omitempty"`
		*Alias
	}{
		ID:          id,
		PublishDate: date,
		Alias:       (*Alias)(article),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	article.ID = bson.ObjectIdHex(aux.ID)

	if aux.PublishDate != "" {
		date := strings.Split(aux.PublishDate, "T")
		t, err := time.Parse("2006-01-02", date[0])
		if err != nil {
			return err
		}

		article.PublishDate = t
	}

	return nil
}

//UnmarshalJSON custom Unmarshal function for article model
func (article *Article) UnmarshalJSON1(j []byte) error {
	var articleMap map[string]interface{}

	err := json.Unmarshal(j, &articleMap)
	if err != nil {
		return err
	}

	for k, v := range articleMap {
		switch strings.ToLower(k) {
		case "id":
			article.ID = bson.ObjectIdHex(v.(string))
		case "title":
			article.Title = v.(string)
		case "author":
			article.Author = v.(string)
		case "url":
			article.URL = v.(string)
		case "content":
			article.Content = v.(string)
		case "banner":
			article.Banner = v.(string)
		case "datasource":
			article.DataSource = v.(string)
		case "publishdate":
			t, err := time.Parse("2006-01-02", v.(string))
			if err != nil {
				return err
			}

			article.PublishDate = t
		case "categories":
			cats := make([]string, 0)
			for _, v := range v.([]interface{}) {
				cats = append(cats, v.(string))
			}

			article.Categories = cats
		case "tags":
			tags := make([]string, 0)
			for _, v := range v.([]interface{}) {
				tags = append(tags, v.(string))
			}

			article.Tags = tags
		}
	}

	return nil
}
