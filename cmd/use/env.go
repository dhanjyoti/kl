package use

import (
	"fmt"

	"github.com/kloudlite/kl/pkg/ui/fzf"
	"github.com/kloudlite/kl/pkg/ui/text"

	"github.com/kloudlite/kl/domain/apiclient"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/spf13/cobra"
)

/*
steps to be peformed:
1. list envs for current team from apic
2. pick one using fzf
3. persist selected env in session
*/

var switchCmd = &cobra.Command{
	Use:   "env",
	Short: "Switch to a different environment",
	Run: func(cmd *cobra.Command, args []string) {
		if err := switchEnv(cmd, args); err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func switchEnv(*cobra.Command, []string) error {
	apic, err := apiclient.New()
	if err != nil {
		return err
	}

	klFile, err := apic.GetFClient().GetKlFile()
	if err != nil {
		return err
	}

	currentTeam, err := apic.GetFClient().GetDataContext().GetTeam()
	if err != nil {
		return fn.NewE(err)
	}

	envs, err := apic.ListEnvs(currentTeam)
	if err != nil {
		return fn.NewE(err)
	}

	env, err := fzf.FindOne(
		envs,
		func(env apiclient.Env) string {
			displayName := fmt.Sprintf("%-40s", env.DisplayName)
			name := fmt.Sprintf("%-30s", env.Metadata.Name)

			if env.ClusterName == "" {
				name := fmt.Sprintf("%-30s", fmt.Sprintf("%s (template)", env.Metadata.Name))
				return fmt.Sprintf("%s %s", name, displayName)
			}
			return fmt.Sprintf("%s %s", name, displayName)
		},
		fzf.WithPrompt("Select Environment > "),
	)

	if err != nil {
		return fn.NewE(err)
	}

	dctx := apic.GetFClient().GetDataContext()

	wc, err := apic.GetFClient().GetWsContext()
	if err != nil {
		return err
	}

	if err := wc.SetEnv(env.Metadata.Name); err != nil {
		return fn.NewE(err)
	}

	ed, err := apic.GetFClient().GetExtraData()
	if err != nil {
		return err
	}

	dnsHostSuffix := ed.GetDnsHostSuffix()

	if dnsHostSuffix != "" {
		env, err := dctx.GetEnv()
		if err != nil {
			return err
		}
		team, err := dctx.GetTeam()
		if err != nil {
			return err
		}

		dctx.SetSearchDomain(fmt.Sprintf("%s.%s.%s", env, team, dnsHostSuffix))
	}

	if klFile.DefaultEnv == "" {
		klFile.DefaultEnv = env.Metadata.Name
		if err := klFile.Save(); err != nil {
			return err
		}
	}
	fn.Log(text.Bold(text.Green("\nSelected Environment:")),
		text.Blue(fmt.Sprintf("\n%s (%s)", env.DisplayName, env.Metadata.Name)),
	)

	return nil
}

func init() {
	switchCmd.Aliases = append(switchCmd.Aliases, "switch")

	switchCmd.Flags().StringP("envname", "e", "", "environment name")
	switchCmd.Flags().StringP("team", "a", "", "team name")
}
