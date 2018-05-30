package articles

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"

	"fmt"
	"github.com/evcraddock/goarticles/models"
	"gopkg.in/mgo.v2/bson"
)

type Controller struct {
	repository Repository
}

func CreateArticleController(router *mux.Router, config models.Config) {
	server := fmt.Sprintf("%v:%v", config.DatabaseServer, config.DatabasePort)
	repository := CreateArticleRepository(server, config.DatabaseName)
	controller := Controller{repository: *repository}

	router.HandleFunc("/api/articles", controller.GetAll).Methods("GET")
	router.HandleFunc("/api/articles/{id}", controller.Get).Methods("GET")
	router.HandleFunc("/api/articles", controller.Add).Methods("POST")
	router.HandleFunc("/api/articles/{id}", controller.Update).Methods("PUT")
	router.HandleFunc("/api/articles/{id}", controller.Delete).Methods("DELETE")
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	articles := c.repository.GetArticles()

	data, _ := json.Marshal(articles)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	return
}

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
