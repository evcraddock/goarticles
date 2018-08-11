package repo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/evcraddock/goarticles"
)

//ArticleRepository model
type ArticleRepository struct {
	Server       string
	DatabaseName string
}

//CreateArticleRepository creates a new repository
func CreateArticleRepository(server, databaseName string) *ArticleRepository {
	return &ArticleRepository{
		Server:       server,
		DatabaseName: databaseName,
	}
}

//GetArticles returns queried articles from database
func (r *ArticleRepository) GetArticles(query map[string]interface{}) goarticles.Articles {

	session, err := mgo.Dial(r.Server)
	if err != nil {
		log.Warn("Failed to establish connection to database:", err)
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("articles")
	results := goarticles.Articles{}
	if err := c.Find(query).All(&results); err != nil {
		log.Warn("Failed to write results:", err)
	}

	log.Debug("GetArticles returned ", len(results), " articles")

	return results
}

//GetArticle returns article by Id
func (r *ArticleRepository) GetArticle(id string) (*goarticles.Article, error) {
	session, err := mgo.Dial(r.Server)

	if err != nil {
		log.Warn("Failed to establish connection to database:", err)
		return nil, err
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("articles")

	result := goarticles.Article{}
	if err := c.FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		log.Warn("Failed to write results:", err)
		return nil, err
	}

	log.Debug("GetArticle returned ", result.ID)

	return &result, nil
}

//AddArticle add article to database
func (r *ArticleRepository) AddArticle(article goarticles.Article) (*goarticles.Article, error) {
	session, err := mgo.Dial(r.Server)
	defer session.Close()

	article.ID = bson.NewObjectId()
	session.DB(r.DatabaseName).C("articles").Insert(article)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	log.Debug("Added Article ID: ", article.ID)

	return &article, nil
}

//UpdateArticle updates article
func (r ArticleRepository) UpdateArticle(article goarticles.Article) (*goarticles.Article, error) {
	session, err := mgo.Dial(r.Server)
	defer session.Close()
	session.DB(r.DatabaseName).C("articles").UpdateId(article.ID, article)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	log.Debug("Updated Article ID: ", article.ID)

	return &article, nil
}

//DeleteArticle deletes article
func (r ArticleRepository) DeleteArticle(id string) error {
	session, err := mgo.Dial(r.Server)
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("article doesn't exist")
	}

	oid := bson.ObjectIdHex(id)
	if err = session.DB(r.DatabaseName).C("articles").RemoveId(oid); err != nil {
		log.Warn(err)
		return err
	}

	log.Debug("Delete Article ID: ", id)

	return nil
}