package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DirectorURL    string `json:"directorURL"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	CAPath         string `json:"caPath"`
	DeploymentName string `json:"deploymentName"`
	KillInterval   int    `json:"killInterval"`
}

func parseConfig(configPath string) Config {
	c := Config{}
	confJ, _ := ioutil.ReadFile(configPath)
	// TODO: Handle missing file/parse error
	json.Unmarshal(confJ, &c)
	return c
}
