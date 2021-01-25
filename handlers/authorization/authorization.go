package authorization

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"assignment/constants"
	"assignment/db"
)

const (
	blankString    = ""
	ownerPathParam = "owner"
)

//Redirector defines Redirect signature
type Redirector interface {
	Redirect(ctx context.Context) (string, error)
}

//AuthHandler handles authorization requests
type AuthHandler struct {
	uc Redirector
}

//CreateHandler creates authorization handler
func CreateHandler(uc Redirector) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	owner := params[ownerPathParam]
	db.UserMapLock.RLock()
	_, ok := db.UserToToken[owner]
	db.UserMapLock.RUnlock()
	if ok {
		w.Write([]byte("User Authenticated"))
		return
	}
	ctx := context.WithValue(context.Background(), constants.UserIDKey, owner)
	reDirectURL, err := ah.uc.Redirect(ctx)
	if err != nil {
		//TODO
		fmt.Println("Error in getting redirect URL:", err)
	}
	http.Redirect(w, r, reDirectURL, http.StatusMovedPermanently)
}
