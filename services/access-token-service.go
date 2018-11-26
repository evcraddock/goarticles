package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

//AccessTokenService service used to handle interactions with the API
type AccessTokenService struct {
	AuthURL     string
	Auth        AuthRequestBody
	AccessToken string
	Cache       cache.Cache
}

//AuthRequestBody request headers for getting authorization token
type AuthRequestBody struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
}

//AuthResponse response from the authorization service
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

//NewAccessTokenService create new article service
func NewAccessTokenService(authURL string, authRequest AuthRequestBody) *AccessTokenService {
	return &AccessTokenService{
		AuthURL: authURL,
		Auth:    authRequest,
		Cache:   *cache.New(5*time.Minute, 10*time.Minute),
	}
}

//GetAccessToken get access token
func (s *AccessTokenService) GetAccessToken() string {

	if token, found := s.getTokenFromCache(); found {
		return token
	}

	b, _ := json.Marshal(s.Auth)
	res, err := http.Post(s.AuthURL, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Error("error getting token: " + err.Error())
	}

	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	authRes := &AuthResponse{}
	if err := json.Unmarshal(resBody, authRes); err != nil {
		log.Error("error decoding body: " + err.Error())
	}

	s.Cache.Set("token", authRes.AccessToken, cache.DefaultExpiration)

	return authRes.AccessToken
}

func (s *AccessTokenService) getTokenFromCache() (string, bool) {

	token := ""

	x, found := s.Cache.Get("token")
	if found {
		token = x.(string)
	}

	return token, found
}
