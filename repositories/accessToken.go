package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"assignment/constants"
	"assignment/db"
	"assignment/entity"
	"assignment/httprest"
)

const (
	accessTokenURL = "%s/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&state=%s"
)

//CreateGithubTokenImpl creates instance of GithubTokenImpl
func CreateGithubTokenImpl() *GithubTokenImpl {
	return &GithubTokenImpl{}
}

//GithubTokenImpl implements AccessToken()
type GithubTokenImpl struct {
}

//AccessToken gets access token
func (gt *GithubTokenImpl) AccessToken(ctx context.Context, oAuthValues entity.OAuthRedirectValues) error {
	code, ok := ctx.Value(constants.CodeQueryParam).(string)
	if !ok || code == "" {
		return errors.New("Code not found")
	}
	url := fmt.Sprintf(accessTokenURL, oAuthValues.GithubOAuthURL, oAuthValues.ClientID, oAuthValues.ClientSecret,
		code, oAuthValues.RedirectURL, oAuthValues.State)
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	fmt.Println("Url = ", url)
	resp, err := httprest.CallRestAPI(http.MethodPost, url, headers, nil)
	if err != nil {
		return err
	}

	apiOutput, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in reading response ", err)
		return err
	}
	token := &entity.AccessToken{}
	err = json.Unmarshal(apiOutput, token)
	if err != nil {
		return err
	}
	if token.AccessToken == "" {
		return errors.New("Access token not found")
	}

	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok || userID == "" {
		return errors.New("UserID not set")
	}

	db.UserMapLock.Lock()
	db.UserToToken[userID] = token.AccessToken
	db.UserMapLock.Unlock()

	return nil
}
