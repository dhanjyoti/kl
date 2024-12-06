package pkg

import (
	"fmt"

	"github.com/kloudlite/kl/domain/fileclient"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/nixpkghandler"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add new package",
	Run: func(cmd *cobra.Command, args []string) {
		fc, err := fileclient.New()
		if err != nil {
			fn.PrintError(err)
			return
		}

		if err := addPackages(fc, cmd, args); err != nil {
			fn.PrintError(err)
			return
		}

	},
}

func addPackages(fc fileclient.FileClient, cmd *cobra.Command, args []string) error {
	fc, err := fileclient.New()
	if err != nil {
		return fn.NewE(err)
	}

	name := fn.ParseStringFlag(cmd, "name")
	if name == "" && len(args) > 0 {
		name = args[0]
	}

	if name == "" {
		return fn.Error("name is required")
	}

	pc, err := nixpkghandler.New(cmd)
	if err != nil {
		return fn.NewE(err)
	}

	name, hashpkg, err := pc.Find(name)
	if err != nil {
		return fn.NewE(err)
	}

	// download and update lockfile
	if err := pc.AddPackage(name, hashpkg); err != nil {
		return fn.NewE(err)
	}

	fn.Println(fmt.Sprintf("Package %s is added successfully", name))
	return nil
}

func init() {
	addCmd.Flags().StringP("name", "n", "", "name of the package to install")
}
