package ssl

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

type CaProjectConfig struct {
	C            string `json:"country"`
	O            string `json:"organization"`
	ST           string `json:"state"`
	L            string `json:"locality"`
	CaOU         string `json:"caOrganizationUnit"`
	CaCN         string `json:"caCommonName"`
	EmailAddress string `json:"emailAddress"`
}

type Env struct {
	*OpenSSLEnv
	ProjectRoot string
	CaProjectConfig
}

func (env *Env) SwitchToProject(path string) error {
	var pcd = filepath.Join(path, "ssl.config.json")
	if _, err := os.Stat(pcd); err != nil {
		return err
	}
	b, err := os.ReadFile(pcd)
	if err != nil {
		return err
	}
	var newConf CaProjectConfig
	err = json.Unmarshal(b, &newConf)
	if err != nil {
		return err
	}
	env.CaProjectConfig = newConf
	env.ProjectRoot = path
	return nil
}

func CreateEnv() *Env {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &Env{
		OpenSSLEnv: &OpenSSLEnv{
			CmdName: "openssl",
			Out:     bytes.NewBuffer(make([]byte, 0, 1024)),
			Err:     bytes.NewBuffer(make([]byte, 0, 1024)),
			wdStack: []string{wd},
		},
	}
}
