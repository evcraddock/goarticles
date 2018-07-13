package models

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

//UnmarshalJSON custom Unmarshal function for article model
func (article *Article) UnmarshalJSON(j []byte) error {
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
