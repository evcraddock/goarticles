package articles

import (
	"github.com/evcraddock/goarticles/services"
	"github.com/gorilla/mux"
)

//CreateRoutes adds routes for the article controller
func CreateRoutes(r *mux.Router, controller ArticleController, auth services.Authorization) {
	r.HandleFunc("/api/articles", auth.Authorize(controller.GetAll)).Methods("GET")
	r.HandleFunc("/api/articles/{id}", auth.Authorize(controller.GetByID)).Methods("GET")
	r.HandleFunc("/api/articles", auth.Authorize(controller.Add)).Methods("POST")
	r.HandleFunc("/api/articles/{id}", auth.Authorize(controller.Update)).Methods("PUT")
	r.HandleFunc("/api/articles/{id}", auth.Authorize(controller.Delete)).Methods("DELETE")
}
