package cmd

import (
	"fmt"
	"github.com/zcubbs/oauth-showcase/client/configs"
	"github.com/zcubbs/oauth-showcase/client/internal"
	"html/template"
	"log"
	"net/http"
)

type UserDetails struct {
	Username string
	Password string
}

var (
	tmpl *template.Template
)

func Start() {
	log.Println("Starting client...")
	configs.PrintConfig()
	// load templates
	tmpl = template.Must(template.ParseFiles("static/index.html"))

	http.HandleFunc("/", clientLoginHandler)

	log.Println(fmt.Sprintf("Client is running at %d port", configs.Cfg.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configs.Cfg.Port), nil))
}

func clientLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Println("error executing template", err)
			return
		}
		return
	}

	details := UserDetails{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	// fetch token for user
	token, err := internal.PerformPasswordGrant(details.Username, details.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var endpoints []map[string]string

	// call secure endpoints
	for _, endpoint := range configs.Cfg.SecureEndpoints {
		response, err := internal.CallSecureEndpoint(
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
