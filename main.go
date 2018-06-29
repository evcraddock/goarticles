package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
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
	r := mux.NewRouter()
	r.StrictSlash(true)
	setupRoutes(r, config)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.Server.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(r),
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
	auth := services.NewAuthorization(config)

	articleController := articles.CreateArticleController(*config)

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	r.HandleFunc("/api/articles", auth.Authorize(articleController.GetAll)).Methods("GET")
	r.HandleFunc("/api/articles/{id}", auth.Authorize(articleController.GetByID)).Methods("GET")
	r.HandleFunc("/api/articles", auth.Authorize(articleController.Add)).Methods("POST")
	r.HandleFunc("/api/articles/{id}", auth.Authorize(articleController.Update)).Methods("PUT")
	r.HandleFunc("/api/articles/{id}", auth.Authorize(articleController.Delete)).Methods("DELETE")

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
