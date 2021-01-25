package authorization

import (
	"context"
	"fmt"
	"net/http"

	"assignment/db"
	"assignment/constants"
	
)

//Authorizer defines Authorize signature
type Authorizer interface {
	Authorize(ctx context.Context) error
}

func CreateAccessTokenHandler(uc Authorizer) *AccessTokenHandler {
	return &AccessTokenHandler{
		uc:uc,
	}
}

//AccessTokenHandler handler function for /callback
type AccessTokenHandler struct {
	uc Authorizer
}

func (ath *AccessTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	code := r.URL.Query().Get(constants.CodeQueryParam)
	state := r.URL.Query().Get(constants.StateQueryParam)
	db.StateMapLock.RLock()
	userID,ok := db.StateToUser[state]
	db.StateMapLock.RUnlock()
	if !ok {
		w.Write([]byte("Kindly authenticate first, call : http://localhost:8080/github/{owner}/authorize"))
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx,constants.CodeQueryParam,code)
	ctx = context.WithValue(ctx,constants.StateQueryParam,state)
	ctx = context.WithValue(ctx,constants.UserIDKey,userID)
	
	err := ath.uc.Authorize(ctx)
	if err != nil {
		fmt.Println("Error in Authorize function", err)
		return
	}
	w.Write([]byte("Authorization Success"))
}
