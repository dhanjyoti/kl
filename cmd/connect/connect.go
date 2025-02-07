package connect

import (
	"bufio"
	"github.com/kloudlite/kl/domain/envclient"
	"github.com/kloudlite/kl/k3s"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/spinner"
	"github.com/kloudlite/kl/pkg/ui/text"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var Command = &cobra.Command{
	Use:   "connect",
	Short: "start the wireguard connection",
	Long:  "This command will start the wireguard connection",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := startWg(cmd); err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func startWg(cmd *cobra.Command) error {
	defer spinner.Client.UpdateMessage("connecting your device")()
	k3sClient, err := k3s.NewClient(cmd)
	if err != nil {
		return fn.NewE(err)
	}

	if envclient.InsideBox() {
		if err = fn.ExecNoOutput("wg-quick down kl-vpn 2> /dev/null | echo already down"); err != nil {
			return fn.NewE(err)
		}

		if err = fn.ExecNoOutput("wg-quick up kl-vpn"); err != nil {
			return fn.NewE(err)
		}

		if err = fn.ExecNoOutput("wg-quick down kl-workspace-wg 2> /dev/null | echo already down"); err != nil {
			return fn.NewE(err)
		}

		if err = fn.ExecNoOutput("wg-quick up kl-workspace-wg"); err != nil {
			return fn.NewE(err)
		}

		return nil
	}

	//r, _ := k3sClient.CheckK3sRunningLocally()
	//if !r {
	//	if envclient.InsideBox() {
	//		if err = fn.ExecNoOutput("wg-quick down kl-vpn 2> /dev/null | echo already down"); err != nil {
	//			return fn.NewE(err)
	//		}
	//
	//		if err = fn.ExecNoOutput("wg-quick up kl-vpn"); err != nil {
	//			return fn.NewE(err)
	//		}
	//	}
	//	return nil
	//}

	//if !envclient.InsideBox() {
	if err := k3sClient.RestartWgProxyContainer(); err != nil {
		return fn.NewE(err)
	}
	//return nil
	//}

	//if err = fn.ExecNoOutput("wg-quick down kl-vpn 2> /dev/null | echo already down "); err != nil {
	//	return fn.NewE(err)
	//}
	//
	//if err = fn.ExecNoOutput("wg-quick up kl-vpn"); err != nil {
	//	return fn.NewE(err)
	//}

	//fc, err := fileclient.New()
	//if err != nil {
	//	return err
	//}
	//k3sTracker, err := fc.GetK3sTracker()
	//if err == nil {
	//	if ChekcWireguardConnection() && k3sTracker.WgConnection {
	//		return nil
	//	}
	//}

	//if err = fn.ExecNoOutput("wg-quick down kl-workspace-wg 2> /dev/null | echo already down"); err != nil {
	//	return fn.NewE(err)
	//}
	//
	//if err := k3sClient.RestartWgProxyContainer(); err != nil {
	//	return fn.NewE(err)
	//}
	//
	//if err = fn.ExecNoOutput("wg-quick up kl-workspace-wg"); err != nil {
	//	return fn.NewE(err)
	//}

	fn.Log(text.Green("device connected"))

	return nil
}

func ChekcWireguardConnection() bool {
	file, err := os.Open("/kl-tmp/online.status")
	if err != nil {
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	status, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false
	}

	if strings.TrimSpace(status) == "online" {
		return true
	}

	return false
}
