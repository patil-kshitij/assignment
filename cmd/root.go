package cmd

import (
	"assignment/config"
	"assignment/router"
	"flag"
	"fmt"
	"net/http"
)

const (
	defaultConfigFile = "config/config.json"
)

//RootCmd is default function which should be called when execution begins
func RootCmd() error {
	configFile := flag.String("config", defaultConfigFile, "provide path to json configuration file")
	flag.Parse()
	err := config.ReadConfig(*configFile)
	if err != nil {
		fmt.Println("error in reading config file = ", err)
	}
	r := router.NewRouter()
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error :", err)
	}
	return nil
}
