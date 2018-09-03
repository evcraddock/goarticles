package repo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/evcraddock/goarticles"
	"github.com/evcraddock/goarticles/services"
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
func (r *ArticleRepository) GetArticles(query map[string]interface{}) (*goarticles.Articles, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("articles")
	results := goarticles.Articles{}
	if err := services.NewError(
		c.Find(query).All(&results),
		"error retrieving data",
		"DatabaseError",
		false); err != nil {
		return nil, err
	}

	return &results, nil
}

//GetArticle returns article by Id
func (r *ArticleRepository) GetArticle(id string) (*goarticles.Article, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("articles")
	result := goarticles.Article{}
	if err := services.NewError(c.FindId(
		bson.ObjectIdHex(id)).One(&result),
		"article doesn't exist",
		"NotFound",
		false); err != nil {
		return nil, err
	}

	return &result, nil
}

//AddArticle add article to database
func (r *ArticleRepository) AddArticle(article goarticles.Article) (*goarticles.Article, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()

	article.ID = bson.NewObjectId()
	if err := services.NewError(
		session.DB(r.DatabaseName).C("articles").Insert(article),
		"failed to create article",
		"DatabaseError",
		false); err != nil {
		return nil, err
	}

	log.Debug("Added Article ID: ", article.ID)

	return &article, nil
}

//UpdateArticle updates article
func (r *ArticleRepository) UpdateArticle(article goarticles.Article) (*goarticles.Article, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()
	c := session.DB(r.DatabaseName).C("articles")
	oid, err := r.articleExists(c, article.ID.Hex())
	if err != nil {
		return nil, err
	}

	if err := services.NewError(
		c.UpdateId(oid, article),
		"failed to update article",
		"DatabaseError",
		false); err != nil {
		return nil, err
	}

	log.Debug("Updated Article ID: ", article.ID)

	return &article, nil
}

//DeleteArticle deletes article
func (r *ArticleRepository) DeleteArticle(id string) error {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return err
	}

	defer session.Close()
	c := session.DB(r.DatabaseName).C("articles")
	oid, err := r.articleExists(c, id)
	if err != nil {
		return err
	}

	if err = services.NewError(
		c.RemoveId(oid),
		"failed to delete article",
		"DatabaseError",
		false); err != nil {
		return err
	}

	log.Debug("Delete Article ID: ", oid)

	return nil
}

//ArticleExists check to see if artcle exists in database
func (r *ArticleRepository) ArticleExists(id string) (bool, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return false, err
	}

	defer session.Close()
	c := session.DB(r.DatabaseName).C("articles")
	if _, err := r.articleExists(c, id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *ArticleRepository) articleExists(collection *mgo.Collection, id string) (*bson.ObjectId, error) {
	if !bson.IsObjectIdHex(id) {
		err := services.NewError(fmt.Errorf("invalid id"), "can not find record: invalid id", "DatabaseError", false)
		return nil, err
	}

	oid := bson.ObjectIdHex(id)
	count, err := collection.FindId(oid).Count()
	if err != nil {
		return nil, services.NewError(err, "could not find article", "DatabaseError", false)
	}

	if count < 1 {
		return nil, services.NewError(fmt.Errorf("article does not exist"), "article does not exist", "NotFound", false)
	}

	return &oid, nil
}
