package list

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "accounts | projects | devices | configs | secrets | apps",
	Long: `Using this command you can list multiple resources.
`,
}

func init() {
	Cmd.AddCommand(accountsCmd)
	Cmd.AddCommand(clustersCmd)
	Cmd.AddCommand(projectsCmd)
	Cmd.AddCommand(envsCmd)
	Cmd.AddCommand(configsCmd)
	Cmd.AddCommand(secretsCmd)
	Cmd.AddCommand(appsCmd)
}
