package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	//TODO: Add to config
	var timeoutWait = time.Second*15

	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheck)
	//TODO: Add Routes

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		//TODO: configure host and port
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Printf("Service started on %v", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), timeoutWait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("Service shutting down")
	os.Exit(0)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "HEAD" {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(405)
	}
}
