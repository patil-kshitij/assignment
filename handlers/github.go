package handlers

import (
	"assignment/constants"
	"assignment/db"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type GithubAction interface {
	Perform(ctx context.Context, owner, repo, branchName, token string) error
}

func CreateGithubHandler(uc GithubAction) *GithubHandler {
	return &GithubHandler{
		uc: uc,
	}
}

type GithubHandler struct {
	uc GithubAction
}

func (gh *GithubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	owner := params[constants.OwnerKey]
	if owner == "" {
		w.Write([]byte("Kindly provide owner name"))
		return
	}
	fmt.Println("owner = ", owner)
	db.UserMapLock.RLock()
	token, ok := db.UserToToken[owner]
	db.UserMapLock.RUnlock()
	if !ok {
		w.Write([]byte("Kindly authenticate first, call : http://localhost:8080/github/{owner}/authorize"))
		return
	}

	repo := params[constants.RepoKey]
	if repo == "" {
		w.Write([]byte("Kindly provide repository name"))
		return
	}

	branchName := params[constants.BranchNameKey]
	if branchName == "" {
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
