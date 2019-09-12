# goarticles

> simple command line content management system

## CLI
### Installing

```
go get github.com/evcraddock/goarticles

```

### Configuration
The cli requires the following environment variables be set. The variables are used to identify the goarticles-api server
as well as the security api and credentials for getting a token to be passed to the goarticles-api for each request.


```
CLI_API_URL: {http://localhost:8080}
CLI_AUTH_URL: {https://yourhost.auth0.com/oauth/token}
CLI_GRANT_TYPE: {client_credentials}
CLI_CLIENT_ID: {your-client-id}
CLI_CLIENT_SECRET: {your-client-secret}
CLI_AUTH_AUDIENCE: {https://api.yourdomain.com}
```

### Usage
```
goarticles import -files={folder with front matter files}
```

goarticles will attempt to process every file with the ext .md in the folder specified and all of it's subfolders.
If no files are specified from the command prompt, you will be prompted for the file or folder name.

Files should be in the following format:

```
id: {optional}
title: Article Title
url: url-slug
images:
- image-in-thefolder.jpg
banner: image-in-thefolder.jpg
publishDate: 01/02/2016
author: Your Name
categories:
- categoryname1
- categoryname2
tags:
- tagname1
- tagname2
---
Conent of your article in markdown format
```
* If the id field is specified and there is a record in the database with that id, the record will be updated.
Otherwise a new record will be created
* goarticles assumes that any images are located in the same folder as the markdown file
* only images in the 'images' collection will be uploaded. The banner value should refer to an image in the images collection

## API
#### Installing

```
go build -o $GOPATH/bin/goarticles-api cmd/goarticles-api/goarticles-api.go
```
### Configuration
The api requires the following environment variables be set

##### Server
The api uses MongoDb for storing the article data.

```
GOA_SERVER_PORT: {8080}
GOA_LOG_LEVEL: {info,debug,error}
GOA_DB_ADDRESS: {localhost}
GOA_DB_PORT: {27017}
GOA_DB_DATABASENAME: {articleDB}
GOA_DB_TIMEOUT: {15s}
ORIGIN_ALLOWED: {*}
```

##### Authentication
The api used Auth0.com for authentication. To setup an account follow the instructions for [setting up the client](https://auth0.com/docs/api-auth/config/using-the-auth0-dashboard).
```
GOA_AUTH_DOMAIN: {yourhost.auth0.com}
GOA_AUTH_AUDIENCE: {https://api.yourdomain.com}
```

##### Image Storage
Images are stored using Google Cloud Storage. To setup an account follow the instructions for
[setting up Google Cloud Storage](https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go).

```
GOOGLE_APPLICATION_CREDENTIALS: {/app/gcp.json}
GOA_GCP_PROJECTID: {your-project-id}
GOA_GCP_BUCKETNAME: {articles}

```

### Docker
* The Dockerfile expects your GOOGLE_APPLICATION_CREDENTIALS to be located in the root folder as gcp.json
* Port 8080 is used by default
* A docker compose file is located in the deployments folder
