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
	state, err := getState()
	if err != nil {
		return OAuthValues, err
	}
	clientID, ok := getClientID()
	if !ok {
		config.AppLogger.ErrorLogger.Println("client_id env var not set")
		return OAuthValues, errors.New("ClientID not set")
	}
	OAuthValues.ClientID = clientID
	OAuthValues.GithubOAuthURL = config.AppConfig.GithubOAuthURL
	OAuthValues.RedirectURL = config.AppConfig.OAuthRedirectURL
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		config.AppLogger.ErrorLogger.Printf("%s key is expected, but not present in context", constants.UserIDKey)
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

func getClientID() (string, bool) {
	return os.LookupEnv(constants.ClientIDEnvKey)
}

func getClientSecret() (string, bool) {
	return os.LookupEnv(constants.ClientSecretEnvKey)
}

func getState() (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		config.AppLogger.ErrorLogger.Println("error occurred in getting UUID :", err)
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
	clientID, ok := getClientID()
	if !ok {
		config.AppLogger.ErrorLogger.Println("client_id env var not set")
		return OAuthValues, errors.New("ClientID not set")
	}
	OAuthValues.ClientID = clientID
	clientSecret, ok := getClientSecret()
	if !ok {
		config.AppLogger.ErrorLogger.Println("client_secret env var not set")
		return OAuthValues, errors.New("Client_Secret not set")
	}
	OAuthValues.ClientSecret = clientSecret
	OAuthValues.RedirectURL = config.AppConfig.OAuthRedirectURL

	OAuthValues.State, ok = ctx.Value(constants.StateQueryParam).(string)
	if !ok {
		config.AppLogger.ErrorLogger.Printf("%s key is expected, but not present in context", constants.StateQueryParam)
		return OAuthValues, errors.New(constants.StateQueryParam + " not present in context")
	}

	return OAuthValues, nil
}
