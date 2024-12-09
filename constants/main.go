package constants

import (
	"fmt"
	"os"

	"github.com/kloudlite/kl/flags"

	"github.com/kloudlite/kl/domain/fileclient"
)

const (
	RuntimeLinux   = "linux"
	RuntimeDarwin  = "darwin"
	RuntimeWindows = "windows"
	SocatImage     = "ghcr.io/kloudlite/hub/socat:latest"

	KLDNS                       = "100.64.0.1"
	InterceptWorkspaceServiceIp = "172.18.0.3"
	K3sServerIp                 = "172.18.0.2"
)

// depricated
func GetWireguardImageName() string {
	return fmt.Sprintf("%s/box/wireguard:%s", flags.ImageBase, flags.Version)
}

func GetK3SImageName() string {
	if s := os.Getenv("KL_K3S_IMAGE_NAME"); s != "" {
		return s
	}

	return fmt.Sprintf("%s/k3s:%s", flags.ImageBase, flags.Version)
}

func GetK3sTrackerImageName() string {
	if s := os.Getenv("KL_K3S_TRACKER_IMAGE_NAME"); s != "" {
		return s
	}

	return fmt.Sprintf("%s/k3s-tracker:%s", flags.ImageBase, flags.Version)
}

func GetBoxImageName() string {
	if s := os.Getenv("KL_BOX_IMAGE_NAME"); s != "" {
		return s
	}
	return fmt.Sprintf("%s/box:%s", flags.ImageBase, flags.Version)
}

var (
	BaseURL = func() string {
		baseUrl := flags.DefaultBaseURL
		if s := os.Getenv("KL_BASE_URL"); s != "" {
			return s
		}

		fc, err := fileclient.New()
		if err != nil {
			return baseUrl
		}

		s, err := fc.GetBaseURL()
		if err != nil {
			return baseUrl
		}

		if s != "" {
			return s
		}

		return baseUrl
	}()

	LoginUrl = func() string {
		return fmt.Sprintf("%s/cli-login", BaseURL)
	}()
	ServerURL = func() string {
		return fmt.Sprintf("%s/api/", BaseURL)
	}()

	UpdateURL = func() string {
		return "https://kl.kloudlite.io/kloudlite/kloudlite"
	}()
)

var (
	CoreCliName = "kl"
)

// Temporary Constants
const (
	NetworkService     = "kloudlite-network"
	LocalSearchDomains = "local-search-domains"
)

// new constants kl-v2
const (
	KL_CONTAINER_MARK = "kl-container-mark"
)
