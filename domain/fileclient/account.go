package fileclient

import (
	"errors"
	"os"
	"path"

	fn "github.com/kloudlite/kl/pkg/functions"
	"sigs.k8s.io/yaml"
)

type AccountContext struct {
	ActiveEnv string `json:"activeEnv"`
}

func (a *AccountContext) GetActiveEnv() (string, error) {
	if a.ActiveEnv == "" {
		return "", fn.Errorf("no environment is active")
	}

	return a.ActiveEnv, nil
}

func (a *AccountContext) Save() error {
	confPath, err := getActiveAccountConfigPath()
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(a)
	if err != nil {
		return err
	}

	return writeOnUserScope(path.Join(confPath, "config.yml"), out)
}

func getActiveAccountConfigPath() (string, error) {
	sd, err := getSessionData()
	if err != nil {
		return "", err
	}

	s, err := sd.GetTeam()
	if err != nil {
		return "", err
	}

	configFolder, err := GetConfigFolder()
	if err != nil {
		return "", err
	}

	return path.Join(configFolder, s, "config.yml"), nil
}

func GetActiveAccountContext() (*AccountContext, error) {
	confPath, err := getActiveAccountConfigPath()
	if err != nil {
		return nil, err
	}
	accContext := AccountContext{}

	b, err := ReadFile(confPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := accContext.Save(); err != nil {
				return nil, err
			}

			return &accContext, nil
		}
	}

	if err := yaml.Unmarshal(b, &accContext); err != nil {
		return nil, err
	}

	return &accContext, nil
}

func (f *fclient) CurrentTeamName() (string, error) {
	sd, err := getSessionData()
	if err != nil {
		return "", fn.NewE(err)
	}

	return sd.GetTeam()
}
