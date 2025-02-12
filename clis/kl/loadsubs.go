package kl

import (
	"os/exec"

	"github.com/kloudlite/kl/cmd/auth"
	"github.com/kloudlite/kl/cmd/box"
	"github.com/kloudlite/kl/cmd/clone"
	"github.com/kloudlite/kl/cmd/cluster"
	"github.com/kloudlite/kl/cmd/connect"
	"github.com/kloudlite/kl/cmd/env"
	"github.com/kloudlite/kl/cmd/expose"
	"github.com/kloudlite/kl/cmd/get"
	"github.com/kloudlite/kl/cmd/intercept"
	"github.com/kloudlite/kl/cmd/kubectl"
	"github.com/kloudlite/kl/cmd/list"
	"github.com/kloudlite/kl/cmd/packages"
	"github.com/kloudlite/kl/cmd/runner"
	"github.com/kloudlite/kl/cmd/runner/add"
	set_base_url "github.com/kloudlite/kl/cmd/set-base-url"
	"github.com/kloudlite/kl/cmd/status"
	"github.com/kloudlite/kl/cmd/use"
	"github.com/kloudlite/kl/flags"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	if flags.IsDev() {
		rootCmd.AddCommand(DocsCmd)
	}

	rootCmd.AddCommand(UpdateCmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(box.BoxCmd)

	rootCmd.AddCommand(use.Cmd)
	rootCmd.AddCommand(clone.Cmd)
	rootCmd.AddCommand(env.Cmd)
	rootCmd.AddCommand(runner.InitCommand)
	rootCmd.AddCommand(runner.AttachCommand)
	rootCmd.AddCommand(set_base_url.Cmd)

	rootCmd.AddCommand(intercept.Cmd)

	rootCmd.AddCommand(cluster.Cmd)
	rootCmd.AddCommand(expose.Cmd)

	rootCmd.AddCommand(add.Command)
	rootCmd.AddCommand(status.Cmd)
	rootCmd.AddCommand(packages.Cmd)
	rootCmd.AddCommand(packages.LibCmd)

	rootCmd.AddCommand(connect.Command)

	if flags.IsDev() {
		if _, err := exec.LookPath("k9s"); err == nil {
			rootCmd.AddCommand(kubectl.K9sCmd)
		}

		if _, err := exec.LookPath("kubectl"); err == nil {
			rootCmd.AddCommand(kubectl.KubectlCmd)
		}
	}
}
