package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/evcraddock/goarticles/api/articles"
	"github.com/evcraddock/goarticles/api/health"
	"github.com/evcraddock/goarticles/models"
	"github.com/evcraddock/goarticles/services"
)

func init() {
	config := models.GetConfig()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel)
}

//
//type Jwks struct {
//	Keys []JSONWebKeys `json:"keys"`
//}
//
//type JSONWebKeys struct {
//	Kty string   `json:"kty"`
//	Kid string   `json:"kid"`
//	Use string   `json:"use"`
//	N   string   `json:"n"`
//	E   string   `json:"e"`
//	X5c []string `json:"x5c"`
//}

func main() {
	config := models.GetConfig()
	auth := services.NewAuthorization(config)

	r := mux.NewRouter().StrictSlash(true)

	setupRoutes(r, config)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", config.ServerAddress, config.ServerPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      auth.Middleware.Handler(r),
	}

	go func() {
		log.Info("Service started on ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Info(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), config.TimeoutWait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Info("Service shutting down")
	os.Exit(0)
}

func setupRoutes(r *mux.Router, config *models.Config) {
	articleController := articles.CreateArticleController(*config)
	r.HandleFunc("/api/articles", articleController.GetAll).Methods("GET")
	r.HandleFunc("/api/articles/{id}", articleController.GetByID).Methods("GET")
	r.HandleFunc("/api/articles", articleController.Add).Methods("POST")
	r.HandleFunc("/api/articles/{id}", articleController.Update).Methods("PUT")
	r.HandleFunc("/api/articles/{id}", articleController.Delete).Methods("DELETE")

	health.CreateRoutes(r)
}

//
//func getPemCert(token *jwt.Token, config *models.Config) (string, error) {
//	cert := ""
//	resp, err := http.Get("https://" + config.AuthDomain + "/.well-known/jwks.json")
//
//	if err != nil {
//		log.Debug(err)
//		return cert, err
//	}
//
//	defer resp.Body.Close()
//
//	var jwks = Jwks{}
//	err = json.NewDecoder(resp.Body).Decode(&jwks)
//
//	if err != nil {
//		return cert, err
//	}
//
//	for k := range jwks.Keys {
//		if token.Header["kid"] == jwks.Keys[k].Kid {
//			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
//		}
//	}
//
//	if cert == "" {
//		err := errors.New("unable to find appropriate key")
//		return cert, err
//	}
//
//	return cert, nil
//}
