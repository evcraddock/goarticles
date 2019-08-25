package articles

import (
	"encoding/json"
	"strings"
	"time"

	"fmt"

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

//ValidateRequiredFields check list of required fields
func (article *Article) ValidateRequiredFields() []error {
	errors := make([]error, 0)
	requiredFields := []string{
		"title",
		"author",
		"url",
		"content",
	}

	var mapobj map[string]interface{}
	obj, _ := json.Marshal(article)
	json.Unmarshal(obj, &mapobj)

	for _, v := range requiredFields {
		if _, found := mapobj[v]; found {
			if mapobj[v] != "" {
				continue
			}
		}

		err := fmt.Errorf("%v is required", v)
		errors = append(errors, err)
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

//MarshalJSON custom MarshalJSON for articles
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

//UnmarshalJSON custom UnmarshalJSON for articles
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

	if aux.ID != "" {
		article.ID = bson.ObjectIdHex(aux.ID)
	}

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
