package health

import (
	"github.com/gorilla/mux"
	"net/http"
)

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