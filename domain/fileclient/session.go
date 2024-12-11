package fileclient

import (
	"os"
	"path"

	confighandler "github.com/kloudlite/kl/pkg/config-handler"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/text"
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

	GetSearchDomain() (string, error)
	SetSearchDomain(domain string) error
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

type SessionData struct {
	Session      string `json:"session"`
	Team         string `json:"team,omitempty"`
	Env          string `json:"env,omitempty"`
	SearchDomain string `json:"searchDomain,omitempty"`

	TeamsData map[string]*DeviceData `json:"teamsData,omitempty"`
}

type sed struct {
	*SessionData
	handler confighandler.Config[SessionData]
}

func (c *sed) Clear() error {
	c.SessionData = &SessionData{}
	return c.handler.Write()
}

func (s *sed) GetSearchDomain() (string, error) {
	if s.SearchDomain == "" {
		return "", fn.Errorf("search domain not found")
	}

	return s.SearchDomain, nil
}

func (s *sed) SetSearchDomain(domain string) error {
	s.SearchDomain = domain
	return s.handler.Write()
}

func (s *sed) SetDevice(dev DeviceData) error {
	team, err := s.GetTeam()
	if err != nil {
		return err
	}

	if s.TeamsData == nil {
		s.TeamsData = make(map[string]*DeviceData)
	}

	if s.TeamsData[team] == nil {
		s.TeamsData[team] = &DeviceData{}
	}

	s.TeamsData[team] = &dev
	return s.Save()
}

func (s *sed) GetDevice() (*DeviceData, error) {
	team, err := s.GetTeam()
	if err != nil {
		return nil, err
	}

	if s.TeamsData[team] == nil {
		return nil, fn.Errorf("device not found")
	}

	return s.TeamsData[team], nil
}

func (s *sed) GetWsTeam() (string, error) {
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

func (s *sed) GetTeam() (string, error) {
	if s.Team == "" {
		return "", fn.Errorf("team not found")
	}

	return s.Team, nil
}

func (s *sed) GetSession() (string, error) {
	if s.Session == "" {
		return "", fn.Errorf("session not found")
	}

	return s.Session, nil
}

func (s *sed) GetEnv() (string, error) {
	if s.Env == "" {
		return "", fn.Errorf("env not found, please run `kl use env` to select an environment")
	}
	return s.Env, nil
}

func (s *sed) SetEnv(env string) error {
	s.Env = env
	return s.Save()
}

func (s *sed) SetTeam(team string) error {
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

func (s *sed) SetSession(sess string) error {
	s.Session = sess
	return s.Save()
}

func (s *sed) Save() error {
	return s.handler.Write()
}
