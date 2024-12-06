package fileclient

import "github.com/kloudlite/kl/pkg/functions"

type fclient struct {
	configPath string
}

type FileClient interface {
	GetHostWgConfig() (string, error)
	GetWGConfig() (*WGConfig, error)
	SetWGConfig(config string) error
	CurrentTeamName() (string, error)
	Logout() error
	GetK3sTracker() (*K3sTracker, error)
	GetVpnTeamConfig(team string) (*TeamVpnConfig, error)
	SetVpnTeamConfig(team string, config *TeamVpnConfig) error
	GetClusterConfig(team string) (*TeamClusterConfig, error)
	SetClusterConfig(team string, accClusterConfig *TeamClusterConfig) error
	DeleteClusterData(team string) error
	GetDevice() (*DeviceContext, error)
	SetDevice(device *DeviceContext) error

	GetKlFile() (*KLFileType, error)

	GetLockfile() (*Lockfile, error)

	SelectEnv(ev string) error
	CurrentEnv() (string, error)
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
