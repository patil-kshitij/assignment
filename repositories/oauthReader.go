package repositories

import (
	"assignment/constants"
	"assignment/entity"
	"context"
	"errors"
	"os"

	"assignment/config"
	"assignment/db"

	"github.com/google/uuid"
)



//CreateOAuthEnvConfigReader creates OAuthEnvConfigReader
func CreateOAuthEnvConfigReader() *OAuthEnvConfigReader {
	return &OAuthEnvConfigReader{}
}

//OAuthEnvConfigReader implementation for Get()
type OAuthEnvConfigReader struct{}

//Get value from env and config
func (oar *OAuthEnvConfigReader) Get(ctx context.Context) (entity.OAuthRedirectValues, error) {
	OAuthValues := entity.OAuthRedirectValues{}
	// clientID := os.Getenv(clientIDKey)
	// if clientID == "" {
	// 	return OAuthValues, errors.New("ClientID not found")
	// }
	state, err := getState()
	if err != nil {
		return OAuthValues, err
	}
	clientID,ok := getClientID()
	if !ok {
		return OAuthValues,errors.New("ClientID not set")
	}
	OAuthValues.ClientID = clientID
	OAuthValues.GithubOAuthURL = config.AppConfig.GithubOAuthURL
	OAuthValues.RedirectURL = config.AppConfig.OAuthRedirectURL
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		return OAuthValues, errors.New("Context doesn't contain userID")
	}
	OAuthValues.Login = userID
	OAuthValues.Scope = config.AppConfig.Scopes[config.FullAccessScope]
	OAuthValues.State = state
	OAuthValues.AllowSignUp = config.AppConfig.AllowSignUp
	db.StateMapLock.Lock()
	db.StateToUser[state] = userID
	db.StateMapLock.Unlock()

	return OAuthValues, nil
}

func getClientID() (string,bool) {
	return os.LookupEnv(constants.ClientIDEnvKey)
}

func getClientSecret() (string,bool) {
	return os.LookupEnv(constants.ClientSecretEnvKey)
}

func getState() (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

//CreateOAuthClientSecretReader creates instance for OAuthClientSecretReader
func CreateOAuthClientSecretReader() *OAuthClientSecretReader {
	return &OAuthClientSecretReader{}
}

//OAuthClientSecretReader will read client secret and does not need state
type OAuthClientSecretReader struct{}

//Get value from env and config
func (oacsr *OAuthClientSecretReader) Get(ctx context.Context) (entity.OAuthRedirectValues, error) {
	OAuthValues := entity.OAuthRedirectValues{}
	OAuthValues.GithubOAuthURL = config.AppConfig.GithubOAuthURL
	clientID,ok := getClientID()
	if !ok {
		return OAuthValues,errors.New("ClientID not set")
	}
	OAuthValues.ClientID = clientID
	clientSecret,ok := getClientSecret()
	if !ok {
		return OAuthValues,errors.New("Client_Secret not set")
	}
	OAuthValues.ClientSecret = clientSecret
	OAuthValues.RedirectURL = config.AppConfig.OAuthRedirectURL

	OAuthValues.State = ctx.Value(constants.StateQueryParam).(string)

	return OAuthValues, nil
}
