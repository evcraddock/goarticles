package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evcraddock/goarticles/internal/configs"
	"github.com/evcraddock/goarticles/internal/services"
	"github.com/evcraddock/goarticles/pkg/links"
)

type LinkManager struct {
	URL         string
	AccessToken string
}

func NewLinkManager(config configs.ClientConfiguration) *LinkManager {

	accessTokenService := services.NewAccessTokenService(config.Auth.URL,
		services.AuthRequestBody{
			GrantType:    config.Auth.GrantType,
			ClientID:     config.Auth.ClientID,
			ClientSecret: config.Auth.ClientSecret,
			Audience:     config.Auth.Audience,
		},
	)

	return &LinkManager{
		URL:         config.URL,
		AccessToken: accessTokenService.GetAccessToken(),
	}
}

func (s *LinkManager) CreateLink(link links.Link) error {
	url := s.URL + "/api/links"

	apiLink, err := s.copyTo(&link)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(apiLink)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 201 {
		fmt.Printf("successfully added link: %v \n", link.Title)
		return nil
	}

	return fmt.Errorf("Error adding link with status: %v", res.Status)
}

func (s *LinkManager) copyTo(link *links.Link) (*links.ApiLink, error) {
	apilink := &links.ApiLink{
		Title:      link.Title,
		URL:        link.URL,
		Banner:     link.Banner,
		Categories: link.Categories,
		Tags:       link.Tags,
	}

	return apilink, nil
}
