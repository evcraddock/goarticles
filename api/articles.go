package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"fmt"
	"net/url"

	"github.com/evcraddock/goarticles"
	"github.com/evcraddock/goarticles/repo"
	"github.com/evcraddock/goarticles/services"
	"gopkg.in/mgo.v2/bson"
)

//ArticleController model
type ArticleController struct {
	repository repo.ArticleRepository
}

//CreateArticleController creates controller and sets routes
func CreateArticleController(dbaddress, dbport, dbname string) ArticleController {
	dbserver := fmt.Sprintf("%v:%v", dbaddress, dbport)
	repository := repo.CreateArticleRepository(dbserver, dbname)
	controller := ArticleController{repository: *repository}

	return controller
}

//GetArticleRoutes return list of routes for articles
func (c *ArticleController) GetArticleRoutes() []Route {
	return []Route{
		{"GET", "/api/articles", true, c.GetAll},
		{"GET", "/api/articles/{id}", true, c.GetByID},
		{"POST", "/api/articles", true, c.Add},
		{"PUT", "/api/articles/{id}", true, c.Update},
		{"DELETE", "/api/articles/{id}", true, c.Delete},
	}
}

//GetByID returns article by article Id
func (c *ArticleController) GetByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	article, err := c.repository.GetArticle(id)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(article)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
	log.Info("Get article by id")

	return nil
}

//GetAll returns all queried articles
func (c *ArticleController) GetAll(w http.ResponseWriter, r *http.Request) error {
	vars := r.URL.Query()
	query := c.createQuery(vars)
	articles, err := c.repository.GetArticles(query)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(articles)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	log.Info("GetAll articles")
	return nil
}

//Add adds new article
func (c *ArticleController) Add(w http.ResponseWriter, r *http.Request) error {
	var article goarticles.Article

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return services.NewError(err, "body is invalid", "FormatError", false)
	}

	defer r.Body.Close()
	if err := services.NewError(
		json.Unmarshal(body, &article),
		"error loading data while updating article",
		"FormatError",
		false); err != nil {
		return err
	}

	newArticle, err := c.repository.AddArticle(article)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(newArticle)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

	return nil
}

//Update updates existing article
func (c *ArticleController) Update(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	var article goarticles.Article
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return services.NewError(err, "body is invalid", "FormatError", false)
	}

	defer r.Body.Close()
	if err := services.NewError(
		json.Unmarshal(body, &article),
		"error loading data while updating article",
		"FormatError",
		false); err != nil {
		return err
	}

	articleIDString := bson.ObjectId(article.ID).Hex()
	if articleIDString != id {
		err := fmt.Errorf("invalid idientifier: %v", id)
		return services.NewError(err, "article id does not match url parameter", "ValidationError", false)
	}

	updatedArticle, err := c.repository.UpdateArticle(article)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(updatedArticle)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	return nil
}

//Delete deletes requested article
func (c *ArticleController) Delete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.DeleteArticle(id); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return nil
}

func (c *ArticleController) createQuery(vars url.Values) bson.M {
	query := make(bson.M)

	for k, v := range vars {
		switch k {
		case "categories":
			query[k] = bson.M{"$in": v}
		case "tags":
			query[k] = bson.M{"$in": v}
		default:
			query[k] = v[0]
		}

	}

	return query
}
