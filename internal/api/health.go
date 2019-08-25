package api

import "net/http"

//GetHealthRoutes returns list of health routes
func GetHealthRoutes() []Route {
	return []Route{
		{"GET", "/health", false, HealthCheck},
	}
}

//HealthCheck returns statusok for GET requests
func HealthCheck(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return nil
}
