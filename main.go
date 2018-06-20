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

	"flag"

	"github.com/evcraddock/goarticles/api/articles"
	"github.com/evcraddock/goarticles/api/health"
	"github.com/evcraddock/goarticles/models"
	"github.com/evcraddock/goarticles/services"
)

var config *models.Configuration

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	configFile := flag.String("configfile", "", "yaml configuration file (optional)")
	flag.Parse()

	var err error

	if *configFile != "" {
		config, err = models.LoadConfig(*configFile)
	} else {
		config, err = models.LoadEnvironmentVariables()
	}

	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	setLogLevel(config.Server.LogLevel)
}

func main() {
	auth := services.NewAuthorization(config)

	r := mux.NewRouter().StrictSlash(true)

	setupRoutes(r, config)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", config.Server.Address, config.Server.Port),
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

	ctx, cancel := context.WithTimeout(context.Background(), config.Database.Timeout)
	defer cancel()

	srv.Shutdown(ctx)
	log.Info("Service shutting down")
	os.Exit(0)
}

func setupRoutes(r *mux.Router, config *models.Configuration) {
	articleController := articles.CreateArticleController(*config)
	r.HandleFunc("/api/articles", articleController.GetAll).Methods("GET")
	r.HandleFunc("/api/articles/{id}", articleController.GetByID).Methods("GET")
	r.HandleFunc("/api/articles", articleController.Add).Methods("POST")
	r.HandleFunc("/api/articles/{id}", articleController.Update).Methods("PUT")
	r.HandleFunc("/api/articles/{id}", articleController.Delete).Methods("DELETE")

	health.CreateRoutes(r)
}

func setLogLevel(logLevel string) log.Level {
	switch logLevel {
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}
