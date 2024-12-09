package fileclient

import "github.com/kloudlite/kl/pkg/functions"

type fclient struct {
	configPath string
}

type FileClient interface {
	SaveBaseURL(url string) error
	GetBaseURL() (string, error)
	GetExtraData() (*ExtraData, error)

	GetHostWgConfig() (string, error)
	GetWGConfig() (*WGConfig, error)
	SetWGConfig(config string) error
	// CurrentTeamName() (string, error)
	Logout() error
	GetK3sTracker() (*K3sTracker, error)
	GetClusterConfig(team string) (*TeamClusterConfig, error)
	SetClusterConfig(team string, accClusterConfig *TeamClusterConfig) error
	DeleteClusterData(team string) error
	GetDevice() (*DeviceData, error)
	SetDevice(device *DeviceData) error
	GetSessionData() (*SessionData, error)

	GetKlFile() (*KLFileType, error)

	GetLockfile() (*Lockfile, error)

	SelectEnv(ev string) error
	CurrentEnv() (string, error)

	GetTeam() (string, error)
	GetWsTeam() (string, error)
}

func New() (FileClient, error) {
	configPath, err := GetConfigFolder()
	if err != nil {
		return nil, functions.NewE(err)
	}

	return &fclient{
		configPath: configPath,
	}, nil
}
