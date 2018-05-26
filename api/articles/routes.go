package articles

import (
	"github.com/gorilla/mux"


)

func CreateRoutes(router *mux.Router) {
	repository := CreateArticleRepository()
	controller := Controller{Repository: *repository}

	router.HandleFunc("/api/articles", controller.GetAll).Methods("GET")
	router.HandleFunc("/api/articles", controller.Add).Methods("POST")
}