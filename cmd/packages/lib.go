package packages

import (
	"github.com/spf13/cobra"
)

var LibCmd = &cobra.Command{
	Use:   "lib",
	Short: "libraries util to manage nix libraries of kl box",
}

func init() {
	LibCmd.Aliases = append(LibCmd.Aliases, "libraries")
	LibCmd.Aliases = append(LibCmd.Aliases, "library")

	LibCmd.AddCommand(listCmd)
	LibCmd.AddCommand(addLibCmd)
	LibCmd.AddCommand(rmLibCmd)
	LibCmd.AddCommand(searchCmd)
}
