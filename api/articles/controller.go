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
)

type Controller struct {
	repository Repository
}

func CreateArticleController(router *mux.Router, config models.Config) {
	server := fmt.Sprintf("%v:%v", config.DatabaseServer, config.DatabasePort)
	repository := CreateArticleRepository(server, config.DatabaseName)
	controller := Controller{repository: *repository}

	router.HandleFunc("/api/articles", controller.GetAll).Methods("GET")
	router.HandleFunc("/api/articles", controller.Add).Methods("POST")
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
