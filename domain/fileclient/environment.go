package fileclient

import (
	"errors"
	"os"
	"path"

	confighandler "github.com/kloudlite/kl/pkg/config-handler"
	fn "github.com/kloudlite/kl/pkg/functions"
)

var NoEnvSelected = fn.Errorf("no selected environment")

type LocalEnv struct {
	ActiveEnv string `json:"activeEnv" yaml:"activeEnv" `
}

func (le *LocalEnv) GetActiveEnv() (string, error) {
	if le.ActiveEnv == "" {
		return "", fn.Errorf("no environment selected")
	}

	return le.ActiveEnv, nil
}

func getEnvConfigPath() (string, error) {
	cf := path.Join(path.Dir(getConfigPath()), ".kl")

	if err := os.MkdirAll(cf, os.ModePerm); err != nil {
		return "", err
	}

	return cf, nil
}

func GetEnvSession() (*LocalEnv, error) {
	s, err := getEnvConfigPath()
	if err != nil {
		return nil, err
	}

	le, err := confighandler.ReadConfig[LocalEnv](path.Join(s, "config.yml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			le := &LocalEnv{}
			if err := le.Save(); err != nil {
				return nil, err
			}
			return le, nil
		}
		return nil, err
	}

	return le, nil
}

func (le *LocalEnv) SetEnv(env string) error {
	le.ActiveEnv = env
	sd, err := getSessionData()

	if err != nil {
		return err
	}
	if err := sd.SetEnv(env); err != nil {
		return err
	}

	return le.Save()
}

func (le *LocalEnv) Save() error {
	envCpath, err := getEnvConfigPath()
	if err != nil {
		return err
	}

	return confighandler.WriteConfig(path.Join(envCpath, "config.yml"), *le, 0o644)
}

func (f *fclient) SelectEnv(env string) error {

	le, err := GetEnvSession()
	if err != nil {
		return err
	}

	return le.SetEnv(env)
}

func (f *fclient) CurrentEnv() (string, error) {
	le, err := GetEnvSession()
	if err != nil {
		return "", err
	}

	return le.GetActiveEnv()
}
