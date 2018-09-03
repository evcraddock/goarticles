package repo

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/evcraddock/goarticles/services"
	log "github.com/sirupsen/logrus"
)

//StorageRepository represents a google cloud storage account
type StorageRepository struct {
	client     *storage.Client
	projectID  string
	bucketName string
}

//CreateNewStorage creates StorageRepository object
func CreateNewStorage(projectName, bucketName string) StorageRepository {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)

	if err != nil {
		panic("Failed to create client: " + err.Error())
	}

	store := StorageRepository{
		client:     storageClient,
		projectID:  projectName,
		bucketName: bucketName,
	}

	store.createBucket(ctx, store.bucketName)

	return store
}

// Creates the new bucket.
func (store *StorageRepository) createBucket(ctx context.Context, bucketName string) *storage.BucketHandle {

	bucket := store.client.Bucket(bucketName)

	if err := bucket.Create(ctx, store.projectID, nil); err != nil {
		log.Infof("Failed to create bucket: %v", err)
	}

	return bucket
}

//AddImage adds image to bucket for an article
func (store *StorageRepository) AddImage(ctx context.Context, image string, file multipart.File) error {
	bucket := store.client.Bucket(store.bucketName)
	imgfile := bucket.Object(image)

	ww := imgfile.NewWriter(ctx)

	if _, err := io.Copy(ww, file); err != nil {
		log.Error(err)
		return services.NewError(err, "unable to write image to bucket", "StorageError", false)
	}

	if err := ww.Close(); err != nil {
		log.Error(err)
		return services.NewError(err, "unable to close writer", "StorageError", false)
	}

	return nil
}

//GetImage get image from storage and return as byte array
func (store *StorageRepository) GetImage(ctx context.Context, image string) ([]byte, error) {
	bucket := store.client.Bucket(store.bucketName)
	imgreader, err := bucket.Object(image).NewReader(ctx)

	if err != nil {
		return nil, services.NewError(err, "could not find image", "NotFound", true)
	}

	defer imgreader.Close()

	imgdata, err := ioutil.ReadAll(imgreader)
	if err != nil {
		return nil, services.NewError(err, "could not open file", "StorageError", false)
	}

	return imgdata, nil
}

//DeleteImage delete requested image
func (store *StorageRepository) DeleteImage(ctx context.Context, image string) error {
	bucket := store.client.Bucket(store.bucketName)
	img := bucket.Object(image)
	if err := img.Delete(ctx); err != nil {
		return services.NewError(err, "unable to delete image", "StorageError", false)
	}

	return nil
}
