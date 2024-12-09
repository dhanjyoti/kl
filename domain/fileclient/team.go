package fileclient

import (
	"errors"
	"os"
	"path"

	fn "github.com/kloudlite/kl/pkg/functions"
	"sigs.k8s.io/yaml"
)

func getSessionData() (*SessionData, error) {

	fc, err := New()
	if err != nil {
		return nil, fn.NewE(err, "failed to create file client")
	}

	file, err := readFile(SessionFileName)
	session := SessionData{
		fc: fc,
	}

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

func getActiveTeamConfigPath() (string, error) {
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

func (f *fclient) GetTeam() (string, error) {
	sd, err := getSessionData()
	if err != nil {
		return "", fn.NewE(err)
	}

	return sd.GetTeam()
}

func (f *fclient) GetWsTeam() (string, error) {
	sd, err := getSessionData()
	if err != nil {
		return "", fn.NewE(err)
	}

	return sd.GetWsTeam()
}
