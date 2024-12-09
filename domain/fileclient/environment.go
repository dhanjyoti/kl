package fileclient

import (
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

func (fc *fclient) GetEnvSession() (*LocalEnv, error) {
	return getEnvSession()
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

	le, err := getEnvSession()
	if err != nil {
		return err
	}

	return le.SetEnv(env)
}

func (f *fclient) CurrentEnv() (string, error) {
	le, err := getEnvSession()
	if err != nil {
		return "", err
	}

	return le.GetActiveEnv()
}
