package fileclient

import "github.com/kloudlite/kl/pkg/functions"

type fclient struct {
	configPath string
	session    Session
}

type FileClient interface {
	SaveBaseURL(url string) error
	GetBaseURL() (string, error)

	GetExtraData() (*ExtraData, error)

	GetHostWgConfig() (string, error)
	GetWGConfig() (*WGConfig, error)
	SetWGConfig(config string) error

	Logout() error

	GetK3sTracker() (*K3sTracker, error)
	GetClusterConfig(team string) (*TeamClusterConfig, error)
	SetClusterConfig(team string, accClusterConfig *TeamClusterConfig) error
	DeleteClusterData(team string) error
	GetDevice() (*DeviceData, error)
	SetDevice(device *DeviceData) error
	GetDataContext() Session

	GetKlFile() (*KLFileType, error)

	GetLockfile() (*Lockfile, error)

	SelectEnv(ev string) error
	CurrentEnv() (string, error)

	GetConfigPath() (string, error)
}

func New() (FileClient, error) {
	configPath, err := GetConfigFolder()
	if err != nil {
		return nil, functions.NewE(err)
	}

	sd, err := getSessionData()
	if err != nil {
		sd = &SessionData{}
	}

	return &fclient{
		configPath: configPath,
		session:    sd,
	}, nil
}
