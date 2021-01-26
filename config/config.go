package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	//FullAccessScope represents full access for a repo
	FullAccessScope = "fullAccess"
)

//AppConfig is in-memory configuration
var AppConfig *configuration

//AppLogger logger used in app
var AppLogger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

type configuration struct {
	OAuthRedirectURL string            `json:"oAuthRedirectURL"`
	LoginPrefix      string            `json:"loginPrefix"`
	Scopes           map[string]string `json:"scopes"`
	GithubOAuthURL   string            `json:"githubOAuthURL"`
	AllowSignUp      string            `json:"allow_signup"`
	GithubAPIURL     string            `json:"githubAPIhost"`
	LogFile          string            `json:"logFile"`
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

//LoadApplication loads config and logger
func LoadApplication(configFilePath string) error {
	err := ReadConfig(configFilePath)
	if err != nil {
		return err
	}
	logFile, err := os.OpenFile(AppConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	AppLogger.InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	AppLogger.ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
