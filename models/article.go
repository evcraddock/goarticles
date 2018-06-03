package models

import (
	"encoding/json"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Article model
type Article struct {
	ID          bson.ObjectId `bson:"_id"`
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	URL         string        `json:"url"`
	Content     string        `json:"content"`
	DataSource  string        `json:"dataSource"`
	PublishDate time.Time     `json:"publishDate"`
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
		case "datasource":
			article.DataSource = v.(string)
		case "publishdate":
			t, err := time.Parse("2006-01-02", v.(string))
			if err != nil {
				return err
			}

			article.PublishDate = t
		}
	}

	return nil
}
