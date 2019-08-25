package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/evcraddock/goarticles/internal/services"
	"github.com/evcraddock/goarticles/pkg/articles"
	"github.com/evcraddock/goarticles/pkg/repos"
)

const maxMemory = 1 * 1024 * 1024

//ImageController model
type ImageController struct {
	storage repos.StorageRepository
}

//CreateImageController creates controller and sets routes
func CreateImageController(projectname, bucketname string) ImageController {
	log.Debugf("CreateImageController started")
	storage := repos.CreateNewStorage(projectname, bucketname)
	controller := ImageController{storage: storage}

	log.Debugf("CreateImageController finished")
	return controller
}

//GetImageRoutes returns a list of images routes
func (c *ImageController) GetImageRoutes() []Route {
	return []Route{
		{"GET", "/api/articles/{id}/images/{filename}", false, c.GetByFilename},
		{"POST", "/api/articles/{id}/images", true, c.Add},
		{"DELETE", "/api/articles/{id}/images/{filename}", true, c.DeleteByFilename},
	}
}

//Add adds new image
func (c *ImageController) Add(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	articleID := vars["id"]

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return services.NewError(err, "invalid multipart format", "FormatError", false)
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()

			image := articles.ArticleImage{
				ArticleID: articleID,
				FileName:  fileHeader.Filename,
				File:      file,
			}

			if err := c.storage.AddImage(context.Background(), image.GetPath(), image.File); err != nil {
				return err
			}

		}
	}

	w.WriteHeader(http.StatusAccepted)

	return nil
}

//GetByFilename get image by filename
func (c *ImageController) GetByFilename(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	articleID := vars["id"]
	filename := vars["filename"]

	image := articles.ArticleImage{
		ArticleID: articleID,
		FileName:  filename,
		File:      nil,
	}

	imagefile, err := c.storage.GetImage(context.Background(), image.GetPath())
	if err != nil {
		return err
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+image.FileName)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write(imagefile)

	return nil
}

//DeleteByFilename delete image by filename
func (c *ImageController) DeleteByFilename(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	articleID := vars["id"]
	filename := vars["filename"]

	image := articles.ArticleImage{
		ArticleID: articleID,
		FileName:  filename,
		File:      nil,
	}

	if err := c.storage.DeleteImage(context.Background(), image.GetPath()); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
