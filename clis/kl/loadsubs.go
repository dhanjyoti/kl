package kl

import (

	// "github.com/kloudlite/kl/cmd/box"

	"runtime"

	"github.com/kloudlite/kl/cmd/auth"
	"github.com/kloudlite/kl/cmd/initp"
	"github.com/kloudlite/kl/cmd/lib"
	"github.com/kloudlite/kl/cmd/packages"
	set_base_url "github.com/kloudlite/kl/cmd/set-base-url"
	"github.com/kloudlite/kl/constants"
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

	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(UpdateCmd)
	rootCmd.AddCommand(set_base_url.Cmd)
	rootCmd.AddCommand(initp.InitCommand)

	if runtime.GOOS == constants.RuntimeWindows {
		return
	}

	rootCmd.AddCommand(packages.Cmd)
	rootCmd.AddCommand(lib.Cmd)

	if runtime.GOOS == constants.RuntimeDarwin {
		return
	}

	// rootCmd.AddCommand(list.Cmd)
	// rootCmd.AddCommand(get.Cmd)
	// rootCmd.AddCommand(box.BoxCmd)

	// rootCmd.AddCommand(use.Cmd)
	// rootCmd.AddCommand(clone.Cmd)
	// rootCmd.AddCommand(env.Cmd)
	// rootCmd.AddCommand(runner.AttachCommand)
	//
	// rootCmd.AddCommand(intercept.Cmd)
	//
	// rootCmd.AddCommand(cluster.Cmd)
	// rootCmd.AddCommand(expose.Cmd)
	//
	// rootCmd.AddCommand(add.Command)
	// rootCmd.AddCommand(status.Cmd)
	// rootCmd.AddCommand(packages.Cmd)
	// rootCmd.AddCommand(packages.LibCmd)
	//
	// rootCmd.AddCommand(connect.Command)
	// rootCmd.AddCommand(v2Shell.Command)
	// rootCmd.AddCommand(v2Add.Command)
	// rootCmd.AddCommand(v2Pkg.Command)
	// rootCmd.AddCommand(v2Lib.Command)

	// if flags.IsDev() {
	// 	if _, err := exec.LookPath("k9s"); err == nil {
	// 		rootCmd.AddCommand(kubectl.K9sCmd)
	// 	}
	//
	// 	if _, err := exec.LookPath("kubectl"); err == nil {
	// 		rootCmd.AddCommand(kubectl.KubectlCmd)
	// 	}
	// }
}
