package fileclient

import (
	"errors"
	"os"
	"path"

	"github.com/kloudlite/kl/pkg/functions"
	fn "github.com/kloudlite/kl/pkg/functions"
	"sigs.k8s.io/yaml"
)

type SessionData struct {
	Session string `json:"session"`
	Team    string `json:"team,omitempty"`
}

func (s *SessionData) GetTeam() (string, error) {
	if s.Team == "" {
		return "", fn.Errorf("team not found")
	}

	return s.Team, nil
}

func (s *SessionData) GetSession() (string, error) {
	if s.Session == "" {
		return "", fn.Errorf("session not found")
	}

	return s.Session, nil
}

func (s *SessionData) SetTeam(team string) error {
	s.Team = team

	if team != "" {
		configFolder, err := GetConfigFolder()
		if err != nil {
			return err
		}

		tempdir := path.Join(configFolder, team)

		if err := os.MkdirAll(tempdir, os.ModePerm); err != nil {
			return nil
		}
	}

	return s.Save()
}

func (s *SessionData) SetSession(sess string) error {
	s.Session = sess
	return s.Save()
}

func GetSessionData() (*SessionData, error) {
	file, err := ReadFile(SessionFileName)
	session := SessionData{}

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			b, err := yaml.Marshal(session)
			if err != nil {
				return nil, functions.NewE(err, "failed to marshal session")
			}

			if err := writeOnUserScope(SessionFileName, b); err != nil {
				return nil, functions.NewE(err, "failed to save session")
			}

			return &session, nil
		}
	}

	if err = yaml.Unmarshal(file, &session); err != nil {
		return nil, functions.NewE(err, "failed to unmarshal session")
	}

	return &session, nil
}

func (s *SessionData) Save() error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return writeOnUserScope(SessionFileName, b)
}
