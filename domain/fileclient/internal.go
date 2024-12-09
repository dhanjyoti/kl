package fileclient

import (
	"errors"
	"os"
	"path"

	confighandler "github.com/kloudlite/kl/pkg/config-handler"
)

func getEnvSession() (*LocalEnv, error) {
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

func getEnvConfigPath() (string, error) {
	cf := path.Join(path.Dir(getConfigPath()), ".kl")

	if err := os.MkdirAll(cf, os.ModePerm); err != nil {
		return "", err
	}

	return cf, nil
}
