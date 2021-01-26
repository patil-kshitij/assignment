package handlers

import (
	"assignment/config"
	"assignment/constants"
	"assignment/db"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

//GithubAction defines interface for github action
type GithubAction interface {
	Perform(ctx context.Context, owner, repo, branchName, token string) error
}

//CreateGithubHandler returns instance of GithubHandler
func CreateGithubHandler(uc GithubAction) *GithubHandler {
	return &GithubHandler{
		uc: uc,
	}
}

//GithubHandler handles github request
type GithubHandler struct {
	uc GithubAction
}

func (gh *GithubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config.AppLogger.InfoLogger.Println("Recieved /github/{owner}/{repo}/{branch} call")
	params := mux.Vars(r)
	owner := params[constants.OwnerKey]
	if owner == "" {
		config.AppLogger.ErrorLogger.Fatal("owner argument in /github/{owner}/{repo}/{branch} is blank or not recieved")
		w.Write([]byte("Kindly provide owner name"))
		return
	}
	config.AppLogger.InfoLogger.Println("Recieved owner =", owner)
	db.UserMapLock.RLock()
	token, ok := db.UserToToken[owner]
	db.UserMapLock.RUnlock()
	if !ok {
		config.AppLogger.ErrorLogger.Fatalf("Owner %s called /github/{owner}/{repo}/{branch} api without authenticating", owner)
		w.Write([]byte("Kindly authenticate first, call : http://localhost:8080/github/{owner}/authorize"))
		return
	}

	repo := params[constants.RepoKey]
	if repo == "" {
		config.AppLogger.ErrorLogger.Fatal("repo argument in /github/{owner}/{repo}/{branch} is blank or not recieved")
		w.Write([]byte("Kindly provide repository name"))
		return
	}

	branchName := params[constants.BranchNameKey]
	if branchName == "" {
		config.AppLogger.ErrorLogger.Fatal("branch argument in /github/{owner}/{repo}/{branch} is blank or not recieved")
		w.Write([]byte("Kindly provide branch name"))
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.FileNameKey, constants.DefaultFileName)

	err := gh.uc.Perform(ctx, owner, repo, branchName, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Pull Request created"))
	return

}
