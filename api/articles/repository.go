package articles

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"

	log "github.com/sirupsen/logrus"
)

type Article struct {
	ID     		bson.ObjectId 	`bson:"_id"`
	Title  		string        	`json:"title"`
	Author 		string			`json:"author"`
	Url	 		string        	`json:"url"`
	Content 	string        	`json:"content"`
	DataSource	string 			`json:"dataSource"`
	PublishDate	time.Time		`json:"publishDate"`
}

type Articles []Article

type Repository struct{
	Server string
	DatabaseName string
}


func CreateArticleRepository() *Repository {
	return &Repository{
		Server: "localhost:27017",
		DatabaseName: "articleDB",
	}
}

func (r *Repository) GetArticles() Articles {
	session, err := mgo.Dial(r.Server)
	if err != nil {
		log.Warn("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("articles")
	results := Articles{}
	if err := c.Find(nil).All(&results); err != nil {
		log.Warn("Failed to write results:", err)
	}

	log.Debug("GetArticles returned ", len(results), " articles")

	return results
}

func (r *Repository) AddArticle(article Article) bool {
	session, err := mgo.Dial(r.Server)
	defer session.Close()

	article.ID = bson.NewObjectId()
	session.DB(r.DatabaseName).C("articles").Insert(article)
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Debug("Added Article ID: ", article.ID)

	return true
}