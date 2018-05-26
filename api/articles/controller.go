package articles

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	Repository Repository
}

//func CreateArticleController(router *mux.Router) *Controller {
//	repository := CreateArticleRepository()
//	controller := Controller{Repository: *repository}
//
//	router.HandleFunc("/api/articles", controller.GetAll).Methods("GET")
//	router.HandleFunc("/api/articles", controller.Add).Methods("POST")
//}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	articles := c.Repository.GetArticles()

	data, _ := json.Marshal(articles)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	return
}

func (c *Controller) Add(w http.ResponseWriter, r *http.Request) {
	var article Article

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

	success := c.Repository.AddArticle(article)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	return
}