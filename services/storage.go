package services

import (
	"context"

	log "github.com/sirupsen/logrus"

	"io"

	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/evcraddock/goarticles/models"
)

//Storage represents a google cloud storage account
type Storage struct {
	client     *storage.Client
	projectID  string
	bucketName string
}

//CreateNewStorage creates Storage object
func CreateNewStorage(config models.Configuration) Storage {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	store := Storage{
		client:     storageClient,
		projectID:  config.Storage.Project,
		bucketName: config.Storage.Bucket,
	}

	store.createBucket(ctx, store.bucketName)

	return store
}

// Creates the new bucket.
func (store *Storage) createBucket(ctx context.Context, bucketName string) *storage.BucketHandle {

	bucket := store.client.Bucket(bucketName)

	if err := bucket.Create(ctx, store.projectID, nil); err != nil {
		log.Infof("Failed to create bucket: %v", err)
	}

	return bucket
}

//AddImage adds image to bucket for an article
func (store *Storage) AddImage(ctx context.Context, image models.ArticleImage) error {
	bucket := store.client.Bucket(store.bucketName)
	imgfile := bucket.Object(image.GetPath())

	ww := imgfile.NewWriter(ctx)

	if _, err := io.Copy(ww, image.File); err != nil {
		log.Error(err)
		return err
	}

	if err := ww.Close(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

//GetImage get image from storage and return as byte array
func (store *Storage) GetImage(ctx context.Context, image models.ArticleImage) ([]byte, error) {
	bucket := store.client.Bucket(store.bucketName)
	imgreader, err := bucket.Object(image.GetPath()).NewReader(ctx)

	if err != nil {
		return nil, err
	}

	defer imgreader.Close()

	imgdata, err := ioutil.ReadAll(imgreader)
	if err != nil {
		return nil, err
	}

	return imgdata, nil
}

//DeleteImage delete requested image
func (store *Storage) DeleteImage(ctx context.Context, image models.ArticleImage) error {
	bucket := store.client.Bucket(store.bucketName)
	img := bucket.Object(image.GetPath())
	if err := img.Delete(ctx); err != nil {
		return err
	}

	return nil
}
