package links

import (
	"gopkg.in/mgo.v2/bson"
)

type Link struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Title      string        `json:"title"`
	URL        string        `json:"url"`
	Banner     string        `json:"banner"`
	Categories []string      `json:"categories"`
	Tags       []string      `json:"tags"`
}

type Links []Link
