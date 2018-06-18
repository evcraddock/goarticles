package articles

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"fmt"
	"net/url"

	"github.com/evcraddock/goarticles/models"
	"gopkg.in/mgo.v2/bson"
)

//Controller model
type Controller struct {
	repository Repository
}

//CreateArticleController creates controller and sets routes
func CreateArticleController(config models.Config) Controller {
	server := fmt.Sprintf("%v:%v", config.DatabaseServer, config.DatabasePort)
	repository := CreateArticleRepository(server, config.DatabaseName)
	controller := Controller{repository: *repository}

	return controller
}

//GetByID returns article by article Id
func (c *Controller) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	article, err := c.repository.GetArticle(id)

	if err != nil {
		//TODO: return an error message
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(article)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//GetAll returns all queried articles
func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()
	query := c.createQuery(vars)
	articles := c.repository.GetArticles(query)

	data, _ := json.Marshal(articles)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	return
}

//Add adds new article
func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var article models.Article

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error adding article", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error adding article", err)
	}

	if err := json.Unmarshal(body, &article); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling article data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	newArticle, err := c.repository.AddArticle(article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(newArticle)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

	return
}

//Update updates existing article
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var article models.Article
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Fatalln("Error Updating article", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error updating article", err)
	}

	if err := json.Unmarshal(body, &article); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error unmarshalling data while updating article", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	articleIDString := bson.ObjectId(article.ID).Hex()
	if articleIDString != id {
		//TODO: return error messages
		w.WriteHeader(http.StatusBadRequest)
		log.Warn("article ID", articleIDString, " does not match ID parameter ", id)
		return
	}

	updatedArticle, err := c.repository.UpdateArticle(article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, _ := json.Marshal(updatedArticle)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	return
}

//Delete deletes requested article
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := c.repository.DeleteArticle(id); err != nil {
		//TODO: return error message
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

func (c *Controller) createQuery(vars url.Values) bson.M {
	query := make(bson.M)

	//TODO: create and check white list, parse dates and id differently
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
