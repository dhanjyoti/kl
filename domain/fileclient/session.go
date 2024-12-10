package fileclient

import (
	"os"
	"path"

	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/text"
	"sigs.k8s.io/yaml"
)

type Session interface {
	GetDevice() (*DeviceData, error)
	SetDevice(dev DeviceData) error

	GetSession() (string, error)
	SetSession(sess string) error

	GetTeam() (string, error)
	GetWsTeam() (string, error)
	SetTeam(team string) error

	GetEnv() (string, error)
	SetEnv(env string) error

	Clear() error
}

func (c *fclient) GetDataContext() Session {
	return c.session
}

type DeviceData struct {
	WGconf     string `json:"wg"`
	IpAddress  string `json:"ip"`
	DeviceName string `json:"device"`
}

type EnvCacheData struct{}

type EnvData struct {
	EnvCache EnvCacheData `json:"envCache"`
}

type TeamData struct {
	Device   *DeviceData        `json:"devices,omitempty"`
	EnvDatas map[string]EnvData `json:"envs,omitempty"`
}

type SessionData struct {
	Session string `json:"session"`
	Team    string `json:"team,omitempty"`
	Env     string `json:"env,omitempty"`

	// Devices      map[string]DeviceData `json:"devices,omitempty"`
	// EnvsData     map[string]string     `json:"envsData,omitempty"`

	SelectedEnvs map[string]string    `json:"selectedEnv,omitempty"`
	TeamsData    map[string]*TeamData `json:"teamsData,omitempty"`
}

func (c *SessionData) Clear() error {
	return nil
}

func (s *SessionData) SetDevice(dev DeviceData) error {
	team, err := s.GetTeam()
	if err != nil {
		return err
	}

	if s.TeamsData == nil {
		s.TeamsData = make(map[string]*TeamData)
	}

	if s.TeamsData[team] == nil {
		s.TeamsData[team] = &TeamData{}
	}

	s.TeamsData[team].Device = &dev
	return s.Save()
}

func (s *SessionData) GetDevice() (*DeviceData, error) {
	team, err := s.GetTeam()
	if err != nil {
		return nil, err
	}

	if s.TeamsData[team] == nil || s.TeamsData[team].Device == nil {
		return nil, fn.Errorf("team not found")
	}

	return s.TeamsData[team].Device, nil
}

func (s *SessionData) GetWsTeam() (string, error) {
	if s.Team == "" {
		return "", fn.Errorf("team not found")
	}

	kt, err := getKlFile()
	if err != nil {
		return "", err
	}

	if kt.TeamName != s.Team {
		return "", fn.Errorf("selected team is not same as current working directory, please change selected team using %s", text.Blue("kl use team"))
	}

	return s.Team, nil
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

	if s.SelectedEnvs == nil {
		s.SelectedEnvs = make(map[string]string)
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	s.SelectedEnvs[dir] = env
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

func (s *SessionData) Save() error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return writeOnUserScope(SessionFileName, b)
}
