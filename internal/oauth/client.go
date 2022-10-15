package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"oauth-showcase/configs"
)

func PerformPasswordGrant(username, password string) (*oauth2.Token, error) {
	// fetch token for user
	token, err := getOauthConfig().PasswordCredentialsToken(
		context.Background(),
		username,
		password,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func CallSecureEndpoint(endpoint string, token *oauth2.Token) (string, error) {
	if token == nil {
		return "", errors.New("token is null")
	}

	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", endpoint, token.AccessToken))
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing response body", err)
		}
	}(resp.Body)
	var response interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	jsonConfig, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonConfig), err
}

func getOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     configs.Cfg.AuthClientID,
		ClientSecret: configs.Cfg.AuthClientSecret,
		Scopes:       configs.Cfg.Scopes,
		RedirectURL:  configs.Cfg.RedirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s%s", configs.Cfg.AuthUrl, "/oauth/authorize"),
			TokenURL: fmt.Sprintf("%s%s", configs.Cfg.AuthUrl, "/oauth/token"),
		},
	}
}
