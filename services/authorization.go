package services

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

//Authorization values for validating token
type Authorization struct {
	domain     string
	audience   string
	middleware *jwtmiddleware.JWTMiddleware
}

//Jwks json web key collection
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

//JSONWebKeys json web key properties
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

//NewAuthorization create new authorization
func NewAuthorization(domain, audience string) Authorization {
	log.Debug("NewAuthorization started")
	auth := Authorization{
		domain:   domain,
		audience: audience,
	}

	auth.middleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: auth.validateToken,
		SigningMethod:       jwt.SigningMethodRS256,
		ErrorHandler:        NotAuthorizedError,
	})

	log.Debug("NewAuthorziation finished")
	return auth
}

//Authorize wraps handler with authorization middleware
func (auth *Authorization) Authorize(handler http.Handler) http.Handler {
	return auth.middleware.Handler(handler)
}

//NotAuthorizedError authorization error handler
func NotAuthorizedError(w http.ResponseWriter, r *http.Request, err string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	apiError := NewError(
		fmt.Errorf("%v", err),
		"unable to validate token",
		"Authorization",
		true)

	errorData, _ := json.Marshal(apiError)
	w.WriteHeader(apiError.Status())
	w.Write(errorData)
}

func (auth *Authorization) validateToken(token *jwt.Token) (interface{}, error) {
	aud := auth.audience
	checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

	if !checkAud {
		return token, errors.New("invalid audience")
	}

	iss := "https://" + auth.domain + "/"
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)

	if !checkIss {
		return token, errors.New("invalid issuer")
	}

	cert, err := auth.getPemCert(token)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

func (auth *Authorization) getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + auth.domain + "/.well-known/jwks.json")

	if err != nil {
		log.Debug(err)
		return cert, err
	}

	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
