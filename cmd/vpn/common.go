package vpn

import (
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"

	"github.com/kloudlite/kl/constants"
	"github.com/kloudlite/kl/domain/apiclient"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/wg_vpn"
	wg_svc "github.com/kloudlite/kl/pkg/wg_vpn/wg_service"
)

const (
	ifName string = "utun2464"
)

func startConfiguration(verbose bool, options ...fn.Option) error {
	apic, err := apiclient.New()
	if err != nil {
		return err
	}

	device, err := apic.EnsureDevice()
	if err != nil {
		return err
	}

	if device.WGconf == "" {
		return errors.New("no wireguard config found, please try again in few seconds")
	}

	configuration, err := base64.StdEncoding.DecodeString(device.WGconf)
	if err != nil {
		return err
	}

	if runtime.GOOS == constants.RuntimeWindows {
		if err := wg_svc.StartVpn(configuration); err != nil {
			return err
		}

		return nil
	}

	if err := wg_vpn.Configure(configuration, ifName, verbose); err != nil {
		return err
	}

	if wg_vpn.IsSystemdReslov() {
		if err := wg_vpn.ExecCmd(fmt.Sprintf("resolvectl domain %s %s", device.DeviceName, func() string {
			// TODO: use real domain
			if true {
				// return fmt.Sprintf("%s.svc.cluster.local", device.EnvironmentName)
				return ""
			}

			return "~."
		}()), false); err != nil {
			return err
		}
	}

	return nil
}
