package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var projectConfig struct {
	C            string `json:"country"`
	O            string `json:"organization"`
	ST           string `json:"state"`
	L            string `json:"locality"`
	CaOU         string `json:"caOrganizationUnit"`
	CaCN         string `json:"caCommonName"`
	EmailAddress string `json:"emailAddress"`
}

func loadProjectConfig(path string) {
	var pcd = filepath.Join(path, "ssl.config.json")
	if _, err := os.Stat(pcd); err != nil {
		panic(err)
	}
	b, err := os.ReadFile(pcd)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &projectConfig)
	if err != nil {
		panic(err)
	}
}
