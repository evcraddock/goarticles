package health

import (
	"net/http"

	"github.com/gorilla/mux"
)

//CreateRoutes create health check route
func CreateRoutes(router *mux.Router) {
	router.HandleFunc("/api/health", healthCheck).Methods("GET")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "HEAD" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
