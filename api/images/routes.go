package images

import (
	"github.com/evcraddock/goarticles/services"
	"github.com/gorilla/mux"
)

//CreateRoutes creates routes for image controller
func CreateRoutes(r *mux.Router, controller ImageController, auth services.Authorization) {
	r.HandleFunc("/api/articles/{id}/images/{filename}", auth.Authorize(controller.GetByFilename)).Methods("GET")
	r.HandleFunc("/api/articles/{id}/images", auth.Authorize(controller.Add)).Methods("POST")
	r.HandleFunc("/api/articles/{id}/images/{filename}", auth.Authorize(controller.DeleteByFilename)).Methods("DELETE")
}
