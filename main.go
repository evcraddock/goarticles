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
)

func init() {
	config := models.GetConfig()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel)
}

func main() {
	config := models.GetConfig()

	r := mux.NewRouter().StrictSlash(true)

	setupRoutes(r, config)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", config.ServerAddress, config.ServerPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Info("Service started on ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Warn(err)
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
	articles.CreateArticleController(r, *config)
	health.CreateRoutes(r)

}
