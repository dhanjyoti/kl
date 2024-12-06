package use

import (
	"github.com/kloudlite/kl/domain/apiclient"
	"github.com/kloudlite/kl/domain/fileclient"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/fzf"
	"github.com/spf13/cobra"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "use team",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := UseTeam(cmd); err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func UseTeam(cmd *cobra.Command) error {
	apic, err := apiclient.New()
	if err != nil {
		return err
	}

	teams, err := apic.ListTeams()
	if err != nil {
		return err
	}

	var selectedTeam *apiclient.Team

	if len(teams) == 0 {
		return fn.Error("no teams found")
	} else if len(teams) == 1 {
		selectedTeam = &teams[0]
	} else {
		selectedTeam, err = fzf.FindOne(teams, func(item apiclient.Team) string {
			return item.Metadata.Name
		}, fzf.WithPrompt("Select team to use >"))
		if err != nil {
			return err
		}
	}

	sd, err := fileclient.GetSessionData()
	if err != nil {
		return err
	}

	if err := sd.SetTeam(selectedTeam.Metadata.Name); err != nil {
		return err
	}

	fn.Log("Selected team is", selectedTeam.Metadata.Name)
	return nil
}
