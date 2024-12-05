package packages

import (
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove installed package",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rmCmd.Flags().StringP("name", "n", "", "name of the package to remove")
}
