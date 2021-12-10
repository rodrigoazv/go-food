package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientID     = "app"
	clientSecret = "..."
)

func main() {

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "http://keycloak:8082/auth/realms/demo")

	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:1323/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "example"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(rw, "State does not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))

		if err != nil {
			http.Error(rw, "Can't change token", http.StatusInternalServerError)
			return
		}

		idToken, ok := oauth2Token.Extra("id_token").(string)

		if !ok {
			http.Error(rw, "Can't pick id token", http.StatusInternalServerError)
			return
		}

		res := struct {
			OAuth2Token *oauth2.Token
			IDToken     string
		}{
			oauth2Token, idToken,
		}

		data, _ := json.MarshalIndent(res, "", " ")
		rw.Write(data)

	})

	log.Fatal(http.ListenAndServe(":1323", nil))

}
