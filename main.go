package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"oauth-showcase/configs"
	"oauth-showcase/internal/oauth"
)

type UserDetails struct {
	Email    string
	Password string
}

func init() {
	configs.Bootstrap()
}

var (
	tmpl *template.Template
)

func main() {
	configs.PrintConfig()
	tmpl = template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", loginHandler)

	log.Println(fmt.Sprintf("Client is running at %d port", configs.Cfg.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configs.Cfg.Port), nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Println("error executing template", err)
			return
		}
		return
	}

	details := UserDetails{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// fetch token for user
	token, err := oauth.PerformPasswordGrant(details.Email, details.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var endpoints []map[string]string

	// call secure endpoints
	for _, endpoint := range configs.Cfg.SecureEndpoints {
		response, err := oauth.CallSecureEndpoint(
			fmt.Sprintf("%s%s", configs.Cfg.AuthUrl, endpoint),
			token,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		endpoints = append(endpoints, map[string]string{
			endpoint: response,
		})
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Success         bool
		AccessToken     string
		RefreshToken    string
		Endpoints       []map[string]string
		DefaultUsername string
		DefaultPassword string
	}{
		true,
		token.AccessToken,
		token.RefreshToken,
		endpoints,
		configs.Cfg.Username,
		configs.Cfg.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
