package use

import (
	"fmt"
	"os"

	"github.com/kloudlite/kl2/pkg/functions"
	"github.com/kloudlite/kl2/pkg/ui/fzf"
	"github.com/kloudlite/kl2/server"
	"github.com/kloudlite/kl2/utils"
	"github.com/kloudlite/kl2/utils/klfile"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "use env",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := selectEnvironment()
			if err != nil {
				functions.PrintError(err)
				return
			}
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return
			}
			err = utils.SetEnvAtPath(cwd, utils.Env(args[0]))
			if err != nil {
				functions.PrintError(err)
				return
			}
		}
	},
}

func selectEnvironment() error {
	klFile, err := klfile.GetKlFile("")
	if err != nil {
		return err
	}
	envs, err := server.ListEnvs(functions.Option{
		Key:   "accountName",
		Value: klFile.AccountName,
	})
	if err != nil {
		return err
	}
	selectedEnv, err := fzf.FindOne(envs, func(item server.Env) string {
		return item.Metadata.Name
	}, fzf.WithPrompt("Select an environment: "))
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = utils.SetEnvAtPath(cwd, utils.Env(selectedEnv.Metadata.Name))
	if err != nil {
		return err
	}
	functions.Log(fmt.Sprintf("switched to %s environment", selectedEnv.Metadata.Name))
	return nil
}
