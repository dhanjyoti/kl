package pkgs

import (
	"fmt"
	"github.com/kloudlite/kl2/pkg/ui/text"
	"github.com/kloudlite/kl2/utils"
	"github.com/kloudlite/kl2/utils/devbox"
	"github.com/kloudlite/kl2/utils/envhash"
	"os"
	"strings"

	fn "github.com/kloudlite/kl2/pkg/functions"
	"github.com/kloudlite/kl2/utils/klfile"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove existing package",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := klfile.GetKlFile("")
		if err != nil {
			fn.PrintError(err)
			return
		}
		for i, p := range config.Packages {
			splits := strings.Split(p, "@")
			if splits[0] == args[0] {
				config.Packages = append(config.Packages[:i], config.Packages[i+1:]...)
				if err := klfile.WriteKLFile(*config); err != nil {
					fn.PrintError(err)
					return
				}
				fn.Println(fmt.Sprintf("Removed %s", args[0]))

				cwd, err := os.Getwd()
				if err != nil {
					fn.PrintError(err)
					return
				}
				env, err := utils.EnvAtPath(cwd)
				if err != nil {
					fn.PrintError(err)
					return
				}

				err = envhash.SyncBoxHash(string(env))
				if err != nil {
					fn.PrintError(err)
					return
				}
				if !(os.Getenv("IN_DEV_BOX") == "true") {
					cwd, err := os.Getwd()
					if err != nil {
						fn.PrintError(err)
						return
					}
					_, err = devbox.ContainerAtPath(cwd)
					if err != nil && err.Error() == devbox.NO_RUNNING_CONTAINERS {
						return
					} else if err != nil {
						fn.PrintError(err)
						return
					}
					fn.Printf(text.Yellow("environments may have been updated. to reflect the changes, do you want to restart the container? [Y/n] "))
					if fn.Confirm("Y", "Y") {
						err = devbox.Stop(cwd)
						if err != nil {
							fn.PrintError(err)
							return
						}
						err = devbox.Start(cwd)
						if err != nil {
							fn.PrintError(err)
							return
						}
					}
				}
				return
			}
		}
		fn.PrintError(fmt.Errorf("package %s not found", args[0]))
	},
}
