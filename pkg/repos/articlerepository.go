package repos

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/evcraddock/goarticles/internal/services"
	"github.com/evcraddock/goarticles/pkg/articles"
)

//ArticleRepository model
type ArticleRepository struct {
	Server       string
	DatabaseName string
}

//CreateArticleRepository creates a new repository
func CreateArticleRepository(server, databaseName string) *ArticleRepository {
	log.Debugf("Database Server: %v", server)
	log.Debugf("Database Name: %v", databaseName)

	return &ArticleRepository{
		Server:       server,
		DatabaseName: databaseName,
	}
}

//GetArticles returns queried articles from database
func (r *ArticleRepository) GetArticles(query map[string]interface{}) (*articles.Articles, error) {
	log.Debugf("Connecting to database %v", r.Server)
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()

	c := session.DB(r.DatabaseName).C("articles")
	results := articles.Articles{}
	if err := services.NewError(
		c.Find(query).Sort("-publishdate").All(&results),
		"error retrieving data",
		"DatabaseError",
		false); err != nil {
		return nil, err
	}

	return &results, nil
}

//GetArticle returns article by Id
func (r *ArticleRepository) GetArticle(id string) (*articles.Article, error) {
	session, err := mgo.Dial(r.Server)
	if err := services.NewError(err, "failed to establish connection to database", "DatabaseConnection", false); err != nil {
		return nil, err
	}

	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err := services.NewError(fmt.Errorf("invalid id"), "can not find record: invalid id", "NotFound", false)
		return nil, err
	}

	c := session.DB(r.DatabaseName).C("articles")
	result := articles.Article{}
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
func (r *ArticleRepository) AddArticle(article articles.Article) (*articles.Article, error) {
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
func (r *ArticleRepository) UpdateArticle(article articles.Article) (*articles.Article, error) {
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

	if err = c.RemoveId(oid); err != nil {
		return services.NewError(err, "failed to delete article", "DatabaseError", false)
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
		err := services.NewError(fmt.Errorf("invalid id"), "can not find record: invalid id", "NotFound", false)
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
