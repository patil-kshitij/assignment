package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	//FullAccessScope represents full access for a repo
	FullAccessScope = "fullAccess"
)

//AppConfig is in-memory configuration
var AppConfig *configuration

type configuration struct {
	OAuthRedirectURL string            `json:"oAuthRedirectURL"`
	LoginPrefix      string            `json:"loginPrefix"`
	Scopes           map[string]string `json:"scopes"`
	GithubOAuthURL   string            `json:"githubOAuthURL"`
	AllowSignUp      string            `json:"allow_signup"`
	GithubAPIURL	 string `json:"githubAPIhost"`
}

//ReadConfig reads configuration from given path
func ReadConfig(configFilePath string) error {
	AppConfig = &configuration{}
	byteArr, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteArr, AppConfig)
	if err != nil {
		return err
	}
	return nil
}
