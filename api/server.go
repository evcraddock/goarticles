package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"encoding/json"

	"github.com/evcraddock/goarticles/configs"
	"github.com/evcraddock/goarticles/services"
)

//Route stores route data
type Route struct {
	Method      string
	Path        string
	RequireAuth bool
	HandlerFunc RouteHandlerFunc
}

//RouteHandler custom route handler
type RouteHandler struct {
	HandlerFunc RouteHandlerFunc
}

//RouteHandlerFunc custom handlerfunc for routes
type RouteHandlerFunc func(w http.ResponseWriter, r *http.Request) error

//AddHandler adds function to RouteHandler
func AddHandler(handlerFunc RouteHandlerFunc) RouteHandler {
	return RouteHandler{
		HandlerFunc: handlerFunc,
	}
}

//ServeHTTP required for Handler interface
func (h RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.HandlerFunc(w, r)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		switch e := err.(type) {
		case services.Error:
			if e.ShouldDisplay() {
				errorData, _ := json.Marshal(e)
				w.WriteHeader(e.Status())
				w.Write(errorData)
				return
			}

			w.WriteHeader(e.Status())
		default:

			//TODO: return better error
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}
}

//NewServer create a new http server
func NewServer(config *configs.Configuration) {
	router := NewRouter(config)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.Server.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Info("Service started on ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Info(err.Error())
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

//NewRouter creates a new router
func NewRouter(config *configs.Configuration) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)
	auth := services.NewAuthorization(config.Authentication.Domain, config.Authentication.Audience)

	var routes []Route

	articleCtrl := CreateArticleController(config.Database.Address, config.Database.Port, config.Database.DatabaseName)
	imageCtrl := CreateImageController(config.Storage.Project, config.Storage.Bucket)

	//TODO: Add Not Found Hander
	//TODO: Add method not allowed handler

	routes = append(routes, articleCtrl.GetArticleRoutes()...)
	routes = append(routes, imageCtrl.GetImageRoutes()...)
	routes = append(routes, GetHealthRoutes()...)

	for _, route := range routes {
		handler := AddHandler(route.HandlerFunc)

		if route.RequireAuth {
			handle := auth.Authorize(handler)
			r.Handle(route.Path, handle).Methods(route.Method)

			continue
		}

		r.Handle(route.Path, handler).Methods(route.Method)
	}

	return handleCORS(r)
}

func handleCORS(router *mux.Router) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	return handlers.CORS(headersOk, originsOk, methodsOk)(router)
}
