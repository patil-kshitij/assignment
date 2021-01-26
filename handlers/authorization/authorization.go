package authorization

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"assignment/config"
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
	config.AppLogger.InfoLogger.Println("recieved /github/{owner}/{repo}/{branch} call")
	params := mux.Vars(r)
	owner := params[ownerPathParam]
	db.UserMapLock.RLock()
	_, ok := db.UserToToken[owner]
	db.UserMapLock.RUnlock()
	if ok {
		config.AppLogger.InfoLogger.Println("recieved redundant authentication request from owner:", owner)
		w.Write([]byte("User Authenticated"))
		return
	}
	ctx := context.WithValue(context.Background(), constants.UserIDKey, owner)
	reDirectURL, err := ah.uc.Redirect(ctx)
	if err != nil {
		w.Write([]byte("Error in getting redirect URL" + err.Error()))
		return
	}
	http.Redirect(w, r, reDirectURL, http.StatusMovedPermanently)
}
