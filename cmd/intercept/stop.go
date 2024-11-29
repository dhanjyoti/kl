package intercept

import (
	"github.com/kloudlite/kl/domain/apiclient"
	"github.com/kloudlite/kl/domain/fileclient"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/fzf"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [app_name]",
	Short: "stop tunneling the traffic to your device",
	Long: `stop intercept service to stop tunnel traffic to your device
Examples:
	# close intercept service
  kl intercept stop [app_name]
	`,

	Run: func(cmd *cobra.Command, _ []string) {

		apic, err := apiclient.New()
		if err != nil {
			fn.PrintError(err)
			return
		}

		fc, err := fileclient.New()
		if err != nil {
			fn.PrintError(err)
			return
		}

		currentAcc, err := fc.CurrentTeamName()
		if err != nil {
			fn.PrintError(err)
			return
		}

		currentEnv, err := fc.CurrentEnv()
		if err != nil {
			fn.PrintError(err)
			return
		}

		apps, err := apic.ListServices(currentAcc, currentEnv.Name)
		if err != nil {
			fn.PrintError(err)
			return
		}

		filteredApps := make([]apiclient.Service, 0)
		for _, app := range apps {
			if app.InterceptStatus.Intercepted {
				filteredApps = append(filteredApps, app)
			}
		}
		if len(filteredApps) == 0 {
			fn.Log("no intercepted apps found")
			return
		}

		appToStop, err := fzf.FindOne(filteredApps, func(item apiclient.Service) string {
			return item.Metadata.Name
		}, fzf.WithPrompt("Select service to stop"))
		if err != nil {
			fn.PrintError(err)
			return
		}

		if err := apic.InterceptService(appToStop, false, nil, currentEnv.Name, []fn.Option{
			fn.MakeOption("appName", appToStop.Metadata.Name),
		}...); err != nil {
			fn.PrintError(err)
			return
		}

		//k3sClient, err := k3s.NewClient()
		//if err != nil {
		//	return
		//}

		//if err = k3sClient.StartAppInterceptService(appToStop.Spec.Intercept.PortMappings, false); err != nil {
		//	return
		//}

		fn.Log("intercepted service stopped successfully")
	},
}

func init() {
	// stopCmd.Flags().StringP("service", "a", "", "service name")

	stopCmd.Aliases = append(stopCmd.Aliases, "close", "end", "leave", "quit", "terminate", "exit", "remove", "disconnect")
}
