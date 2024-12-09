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

func getSessionData() (*SessionData, error) {
	file, err := readFile(SessionFileName)
	session := SessionData{}

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			b, err := yaml.Marshal(session)
			if err != nil {
				return nil, fn.NewE(err, "failed to marshal session")
			}

			if err := writeOnUserScope(SessionFileName, b); err != nil {
				return nil, fn.NewE(err, "failed to save session")
			}

			return &session, nil
		}
	}

	if err = yaml.Unmarshal(file, &session); err != nil {
		return nil, fn.NewE(err, "failed to unmarshal session")
	}

	return &session, nil
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

func (f *fclient) CurrentTeamName() (string, error) {
	sd, err := getSessionData()
	if err != nil {
		return "", fn.NewE(err)
	}

	return sd.GetTeam()
}
