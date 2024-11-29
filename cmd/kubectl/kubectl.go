package kubectl

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kloudlite/kl/flags"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/spf13/cobra"
)

var KubectlCmd = &cobra.Command{
	Use:                "kubectl",
	Short:              "kubectl is a command line tool for controlling Kubernetes clusters",
	Hidden:             !flags.IsDev(),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		kconfPath, err := GetPath(cmd)
		if err != nil {
			fn.PrintError(err)
			return
		}

		c := exec.Command("kubectl", args...)
		c.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", kconfPath))

		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			err = nil
			return
		}

	},
}

func init() {
}
