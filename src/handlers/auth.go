package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"home-counter/src/config"
	"home-counter/src/models"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func saveNewUser(userId string, userName string) {
	var u = models.UserData{}
	err := models.DB.QueryRow(context.Background(), `SELECT * FROM "user" u WHERE u.id = $1`, userId).Scan(&u.Name, &u.Id)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			_, err := models.DB.Exec(context.Background(), `INSERT INTO "user" (id, name) VALUES ($1, $2)`, userId, userName)
			if err != nil {
				panic(err)
			}
		default:
			panic(err)
		}
	}
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.Config.AuthUrl+"/")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID: config.Config.AuthId,
		// fixme add env variables
		ClientSecret: config.Config.AuthSecret,
		RedirectURL:  config.Config.AppUrl + "/user",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

func CallbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := models.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: config.Config.AuthId,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	saveNewUser(profile["sub"].(string), profile["name"].(string))

	// Redirect to logged in page
	http.Redirect(w, r, config.Config.AppUrl+"/user", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := models.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ублюдский костыль, потому что авторы auth0 нехорошие люди
	urlStr := authenticator.Config.AuthCodeURL(state)
	urlStruct, _ := url.Parse(urlStr)
	q, _ := url.ParseQuery(urlStruct.RawQuery)
	fmt.Println(q["redirect_uri"])
	q["redirect_uri"] = []string{config.Config.AppUrl + "/auth/callback"}
	urlStruct.RawQuery = q.Encode()

	http.Redirect(w, r, urlStruct.String(), http.StatusTemporaryRedirect)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	logoutUrl, err := url.Parse(config.Config.AuthUrl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutUrl.Path += "/v2/logout"
	parameters := url.Values{}

	returnTo := config.Config.AppUrl
	parameters.Add("returnTo", returnTo)
	parameters.Add("client_id", config.Config.AuthId)
	logoutUrl.RawQuery = parameters.Encode()

	// TODO clean session!, add another envs!
	session, err := models.Store.Get(r, "auth-session")
	if err != nil {
		log.Printf("Can't get session", err)
		http.Error(w, "Can't get session", 500)
		return
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("Failed to delete session", err)
		return
	}

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
