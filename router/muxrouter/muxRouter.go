package muxrouter

import (
	"assignment/handlers"
	"assignment/handlers/authorization"
	"assignment/repositories"
	"assignment/usecases"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	authorizationHandler = "authorization"
	callback             = "callback"

	combined = "combined"
)

//GetMuxRouter returns instance of mux router
func GetMuxRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/github/{owner}/authorize", getHandler(authorizationHandler)).Methods(http.MethodGet)
	router.Handle("/github/callback", getHandler(callback)).Methods(http.MethodGet)
	router.Handle("/github/{owner}/{repo}/{branch}", getHandler(combined)).Methods(http.MethodGet)

	return router
}

func getHandler(key string) http.Handler {
	fmt.Println(key)
	switch key {

	case authorizationHandler:
		repo := repositories.CreateOAuthEnvConfigReader()
		uc := usecases.CreateAuthRedirect(repo)
		return authorization.CreateHandler(uc)

	case callback:
		accessTokenRepo := repositories.CreateGithubTokenImpl()
		oAuthRepo := repositories.CreateOAuthClientSecretReader()
		uc := usecases.CreateGithubToken(oAuthRepo, accessTokenRepo)
		return authorization.CreateAccessTokenHandler(uc)

	case combined:
		repo := repositories.CreateGithubImpl()
		uc := usecases.CreateCombinedUseCase(repo)
		return handlers.CreateGithubHandler(uc)

	}
	return nil
}
