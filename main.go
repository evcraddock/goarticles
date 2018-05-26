package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/evcraddock/goarticles/api/articles"
	"github.com/evcraddock/goarticles/api/health"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	//TODO: Add to config
	var timeoutWait = time.Second * 15

	r := mux.NewRouter().StrictSlash(true)

	setupRoutes(r)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		//TODO: configure host and port
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

	ctx, cancel := context.WithTimeout(context.Background(), timeoutWait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Info("Service shutting down")
	os.Exit(0)
}

func setupRoutes(r *mux.Router) {
	articles.CreateRoutes(r)
	health.CreateRoutes(r)

}
