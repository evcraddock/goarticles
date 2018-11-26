package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/evcraddock/goarticles/configs"
	"github.com/evcraddock/goarticles/services"
)

//ReverseProxy configuration info
type ReverseProxy struct {
	token            string
	websiteDirectory string
	port             string
	target           *url.URL
	proxy            *httputil.ReverseProxy
}

//NewServer new proxy server
func NewServer(config *configs.ProxyConfiguration) ReverseProxy {
	targetURL, err := url.Parse(config.Server.ForwardAPI)
	if err != nil {
		log.Error("Url is bad: " + err.Error())
		os.Exit(1)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	accessTokenService := services.NewAccessTokenService(config.Auth.URL,
		services.AuthRequestBody{
			GrantType:    config.Auth.GrantType,
			ClientID:     config.Auth.ClientID,
			ClientSecret: config.Auth.ClientSecret,
			Audience:     config.Auth.Audience,
		},
	)

	return ReverseProxy{
		token:            accessTokenService.GetAccessToken(),
		websiteDirectory: config.Server.StaticFiles,
		port:             config.Server.Port,
		proxy:            proxy,
		target:           targetURL,
	}
}

//Start start server
func (s *ReverseProxy) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/api/{rest:..*}", s.handleProxyRequest)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(s.websiteDirectory)))

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", s.port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Info("Service started on ", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Info(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()

	server.Shutdown(ctx)
	log.Info("Service shutting down")
	os.Exit(0)

}

func (s *ReverseProxy) panic(res http.ResponseWriter, r *http.Request, err interface{}) {
	log.Error(r.URL.Path, "an error occurred: ", err)
	res.WriteHeader(http.StatusInternalServerError)
}

func (s *ReverseProxy) handleProxyRequest(res http.ResponseWriter, req *http.Request) {

	log.Debugf("url: %s", req.RequestURI)
	if req.Method != "GET" {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req.Method = "GET"
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Host = s.target.Host
	s.proxy.ServeHTTP(res, req)
}
