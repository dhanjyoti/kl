package packages

import (
	"fmt"
	"strings"

	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/nixpkghandler"
	"github.com/kloudlite/kl/pkg/ui/spinner"
	"github.com/kloudlite/kl/pkg/ui/table"
	"github.com/kloudlite/kl/pkg/ui/text"

	"github.com/spf13/cobra"
)

var PackageSearchCmd = &cobra.Command{
	Use:   "search [name]",
	Short: "search for a package|library",
	Run: func(cmd *cobra.Command, args []string) {
		if err := searchPackages(cmd, args); err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func searchPackages(cmd *cobra.Command, args []string) error {
	name := fn.ParseStringFlag(cmd, "name")
	if name == "" && len(args) > 0 {
		name = args[0]
	}

	if name == "" {
		return fn.Error("name is required")
	}

	defer spinner.Client.UpdateMessage(fmt.Sprintf("searching for package %s", name))()

	pc, err := nixpkghandler.New(cmd)
	if err != nil {
		return fn.NewE(err)
	}

	sr, err := pc.Search(name)
	if err != nil {
		return fn.NewE(err)
	}

	spinner.Client.Pause()

	header := table.Row{table.HeaderText("#"), table.HeaderText("name"), table.HeaderText("versions")}
	rows := make([]table.Row, 0)

	for i, p := range sr.Packages {
		versions := make([]string, 0)
		for j, v := range p.Versions {
			if j >= 10 {
				break
			}

			versions = append(versions, v.Version)
		}

		rows = append(rows, table.Row{
			text.Colored(fmt.Sprint(i+1, "."), 5),
			text.Bold((p.Name)),
			fmt.Sprintf("%s", strings.Join(versions, ", ")),
		})
	}

	fn.Println(table.Table(&header, rows, cmd))

	return nil
}

func init() {
	PackageSearchCmd.Flags().StringP("name", "n", "", "name of the package to remove")

	// TODO: add show-all flag
	// searchCmd.Flags().BoolP("show-all", "a", false, "list all matching packages")
}
