package articles

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"github.com/evcraddock/goarticles/models"
)

type Repository struct {
	Server       string
	DatabaseName string
}

func CreateArticleRepository(server, databaseName string) *Repository {
	return &Repository{
		Server:       server,
		DatabaseName: databaseName,
	}
}

func (r *Repository) GetArticles() models.Articles {
	session, err := mgo.Dial(r.Server)
	if err != nil {
		log.Warn("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("articles")
	results := models.Articles{}
	if err := c.Find(nil).All(&results); err != nil {
		log.Warn("Failed to write results:", err)
	}

	log.Debug("GetArticles returned ", len(results), " articles")

	return results
}

func (r *Repository) AddArticle(article models.Article) (*models.Article, error) {
	session, err := mgo.Dial(r.Server)
	defer session.Close()

	article.ID = bson.NewObjectId()
	session.DB(r.DatabaseName).C("articles").Insert(article)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Debug("Added Article ID: ", article.ID)

	return &article, nil
}
