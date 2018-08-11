package goarticles

import "mime/multipart"

//ArticleImage represents an image for an article
type ArticleImage struct {
	ArticleID string
	FileName  string
	File      multipart.File
}

//GetPath returns bucket path of the file
func (image *ArticleImage) GetPath() string {
	return image.ArticleID + "/" + image.FileName
}
