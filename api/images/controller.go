package images

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/evcraddock/goarticles/models"
	"github.com/evcraddock/goarticles/services"
	log "github.com/sirupsen/logrus"
)

const maxMemory = 1 * 1024 * 1024

//ImageController model
type ImageController struct {
	storage services.Storage
}

//CreateImageController creates controller and sets routes
func CreateImageController(config models.Configuration) ImageController {

	storage := services.CreateNewStorage(config)
	controller := ImageController{storage: storage}

	return controller
}

//Add adds new image
func (c *ImageController) Add(w http.ResponseWriter, r *http.Request) {
	///TODO: check that article exists
	vars := mux.Vars(r)
	articleID := vars["id"]

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()

			image := models.ArticleImage{
				ArticleID: articleID,
				FileName:  fileHeader.Filename,
				File:      file,
			}

			if err := c.storage.AddImage(context.Background(), image); err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

		}
	}

	w.WriteHeader(http.StatusCreated)

	return
}

//GetByFilename get image by filename
func (c *ImageController) GetByFilename(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]
	filename := vars["filename"]

	image := models.ArticleImage{
		ArticleID: articleID,
		FileName:  filename,
		File:      nil,
	}

	imagefile, err := c.storage.GetImage(context.Background(), image)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+image.FileName)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write(imagefile)
}

//DeleteByFilename delete image by filename
func (c *ImageController) DeleteByFilename(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID := vars["id"]
	filename := vars["filename"]

	image := models.ArticleImage{
		ArticleID: articleID,
		FileName:  filename,
		File:      nil,
	}

	if err := c.storage.DeleteImage(context.Background(), image); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
