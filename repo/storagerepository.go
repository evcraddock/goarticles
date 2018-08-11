package repo

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
)

//StorageRepository represents a google cloud storage account
type StorageRepository struct {
	client     *storage.Client
	projectID  string
	bucketName string
}

//CreateNewStorage creates StorageRepository object
func CreateNewStorage(projectname, bucketname string) StorageRepository {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	store := StorageRepository{
		client:     storageClient,
		projectID:  projectname,
		bucketName: bucketname,
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
		return err
	}

	if err := ww.Close(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

//GetImage get image from storage and return as byte array
func (store *StorageRepository) GetImage(ctx context.Context, image string) ([]byte, error) {
	bucket := store.client.Bucket(store.bucketName)
	imgreader, err := bucket.Object(image).NewReader(ctx)

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
func (store *StorageRepository) DeleteImage(ctx context.Context, image string) error {
	bucket := store.client.Bucket(store.bucketName)
	img := bucket.Object(image)
	if err := img.Delete(ctx); err != nil {
		return err
	}

	return nil
}
