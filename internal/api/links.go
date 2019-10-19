package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"

	"github.com/evcraddock/goarticles/internal/services"
	"github.com/evcraddock/goarticles/pkg/links"
	"github.com/evcraddock/goarticles/pkg/repos"
)

type LinkController struct {
	repository repos.LinkRepository
}

func CreateLinkController(dbaddress, dbport, dbname string) LinkController {
	dbserver := fmt.Sprintf("%v:%v", dbaddress, dbport)
	repository := repos.CreateLinkRepository(dbserver, dbname)

	return LinkController{repository: *repository}
}

func (c *LinkController) GetLinkRoutes() []Route {
	return []Route{
		{"GET", "/api/links", false, c.GetAll},
		{"POST", "/api/links", true, c.Add},
		{"DELETE", "/api/links/{id}", true, c.Delete},
	}
}

func (c *LinkController) GetAll(w http.ResponseWriter, r *http.Request) error {
	vars := r.URL.Query()
	query := c.createQuery(vars)
	links, err := c.repository.GetLinks(query)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(links)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	log.Info("GetAll links")
	return nil
}

func (c *LinkController) Add(w http.ResponseWriter, r *http.Request) error {
	var link links.Link

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return services.NewError(err, "body is invalid", "FormatError", false)
	}

	defer r.Body.Close()
	if err := services.NewError(
		json.Unmarshal(body, &link),
		"error loading data while adding link",
		"FormatError",
		false); err != nil {
		return err
	}

	newLink, err := c.repository.AddLink(link)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(newLink)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

	return nil
}

func (c *LinkController) Delete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.repository.DeleteLink(id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	return nil
}

func (c *LinkController) createQuery(vars url.Values) bson.M {
	query := make(bson.M)

	for k, v := range vars {
		switch k {
		case "categories":
			query[k] = bson.M{"$in": v}
		case "tags":
			query[k] = bson.M{"$in": v}
		default:
			query[k] = v[0]
		}

	}

	return query
}
