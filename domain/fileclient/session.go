package fileclient

import (
	"os"
	"path"

	fn "github.com/kloudlite/kl/pkg/functions"
	"sigs.k8s.io/yaml"
)

type DeviceData struct {
	WGconf     string `json:"wg"`
	IpAddress  string `json:"ip"`
	DeviceName string `json:"device"`
}

type SessionData struct {
	Session string                `json:"session"`
	Team    string                `json:"team,omitempty"`
	Env     string                `json:"env,omitempty"`
	Devices map[string]DeviceData `json:"devices,omitempty"`
}

func (s *SessionData) SetDevice(dev DeviceData) error {
	team, err := s.GetTeam()
	if err != nil {
		return err
	}
	if s.Devices == nil {
		s.Devices = make(map[string]DeviceData)
	}

	s.Devices[team] = dev
	return s.Save()
}

func (s *SessionData) GetDevice() (*DeviceData, error) {
	team, err := s.GetTeam()
	if err != nil {
		return nil, err
	}

	if len(s.Devices) == 0 {
		return nil, fn.Errorf("device not found")
	}

	if dev, ok := s.Devices[team]; ok {
		return &dev, nil
	}

	return nil, fn.Errorf("device not found")
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

func (s *SessionData) GetEnv() (string, error) {
	if s.Env == "" {
		return "", fn.Errorf("env not found")
	}
	return s.Env, nil
}

func (s *SessionData) SetEnv(env string) error {
	s.Env = env
	return s.Save()
}

func (s *SessionData) SetTeam(team string) error {
	s.Team = team
	s.Env = ""

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

func (c *fclient) GetSessionData() (*SessionData, error) {
	return getSessionData()
}

func (s *SessionData) Save() error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return writeOnUserScope(SessionFileName, b)
}
