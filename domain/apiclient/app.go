package apiclient

import (
	"fmt"

	"github.com/kloudlite/kl/domain/fileclient"
	"github.com/kloudlite/kl/pkg/functions"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/spinner"
)

var PaginationDefault = map[string]any{
	"orderBy":       "updateTime",
	"sortDirection": "ASC",
	"first":         99999999,
}

type AppSpec struct {
	Services []struct {
		Port int `json:"port"`
	} `json:"services"`
	Intercept *struct {
		Enabled      bool      `json:"enabled"`
		PortMappings []AppPort `json:"portMappings"`
	} `json:"intercept"`
}

type App struct {
	DisplayName string   `json:"displayName"`
	Metadata    Metadata `json:"metadata"`
	Spec        AppSpec  `json:"spec"`
	Status      Status   `json:"status"`
	IsMainApp   bool     `json:"mapp"`
}

type AppPort struct {
	AppPort    int `json:"appPort"`
	DevicePort int `json:"devicePort,omitempty"`
}

func (apic *apiClient) ListApps(teamName string, envName string) ([]App, error) {
	cookie, err := getCookie(fn.MakeOption("teamName", teamName))
	if err != nil {
		return nil, functions.NewE(err)
	}
	respData, err := klFetch("cli_listApps", map[string]any{
		"pq":      PaginationDefault,
		"envName": envName,
	}, &cookie)
	if err != nil {
		return nil, functions.NewE(err)
	}
	if fromResp, err := GetFromRespForEdge[App](respData); err != nil {
		return nil, functions.NewE(err)
	} else {
		return fromResp, nil
	}
}

func (apic *apiClient) InterceptApp(app *App, status bool, ports []AppPort, envName string, options ...fn.Option) error {
	teamName := fn.GetOption(options, "teamName")
	devName := fn.GetOption(options, "deviceName")

	fc, err := fileclient.New()
	if err != nil {
		return functions.NewE(err)
	}

	if teamName == "" {
		kt, err := fc.GetKlFile()
		if err != nil {
			return functions.NewE(err)
		}

		if kt.TeamName == "" {
			return fn.Errorf("team name is required")
		}

		teamName = kt.TeamName
		options = append(options, fn.MakeOption("teamName", teamName))
	}

	if devName == "" {
		sd, err := fc.GetSessionData()
		if err != nil {
			return functions.NewE(err)
		}

		avc, err := sd.GetDevice()
		if err != nil {
			return functions.NewE(err)
		}

		if avc.DeviceName == "" {
			return fmt.Errorf("device name is required")
		}

		devName = avc.DeviceName
	}

	cookie, err := getCookie([]fn.Option{
		fn.MakeOption("teamName", teamName),
	}...)
	if err != nil {
		return functions.NewE(err)
	}

	if len(ports) == 0 {
		if app.Spec.Intercept != nil && len(app.Spec.Intercept.PortMappings) != 0 {
			ports = append(ports, app.Spec.Intercept.PortMappings...)
		} else if len(app.Spec.Services) != 0 {
			for _, v := range app.Spec.Services {
				ports = append(ports, AppPort{
					AppPort:    v.Port,
					DevicePort: v.Port,
				})
			}
		}
	}

	if len(ports) == 0 {
		return fn.Errorf("no ports provided to intercept")
	}

	query := "cli_interceptApp"
	if !app.IsMainApp {
		query = "cli_interceptExternalApp"
	}

	respData, err := klFetch(query, map[string]any{
		"appName":      app.Metadata.Name,
		"envName":      envName,
		"deviceName":   devName,
		"intercept":    status,
		"portMappings": ports,
	}, &cookie)
	if err != nil {
		return functions.NewE(err)
	}

	if _, err := GetFromResp[bool](respData); err != nil {
		return functions.NewE(err)
	} else {
		return nil
	}
}

func (apic *apiClient) RemoveAllIntercepts(options ...fn.Option) error {
	defer spinner.Client.UpdateMessage("Cleaning up intercepts...")()
	devName := fn.GetOption(options, "deviceName")
	teamName := fn.GetOption(options, "teamName")
	currentEnv, err := apic.EnsureEnv()
	if err != nil {
		return functions.NewE(err)
	}

	fc, err := fileclient.New()
	if err != nil {
		return functions.NewE(err)
	}

	if teamName == "" {
		kt, err := fc.GetKlFile()
		if err != nil {
			return functions.NewE(err)
		}

		if kt.TeamName == "" {
			return fn.Errorf("team name is required")
		}

		teamName = kt.TeamName
		options = append(options, fn.MakeOption("teamName", teamName))
	}

	//config, err := apic.fc.GetClusterConfig(teamName)
	//if err != nil {
	//	return functions.NewE(err)
	//}

	if devName == "" {
		sd, err := fc.GetSessionData()
		if err != nil {
			return functions.NewE(err)
		}

		avc, err := sd.GetDevice()
		if err != nil {
			return functions.NewE(err)
		}

		if avc.DeviceName == "" {
			return fn.Errorf("device name is required")
		}

		devName = avc.DeviceName
	}

	cookie, err := getCookie([]fn.Option{
		fn.MakeOption("teamName", teamName),
	}...)
	if err != nil {
		return functions.NewE(err)
	}
	query := "cli_removeDeviceIntercepts"

	respData, err := klFetch(query, map[string]any{
		"envName":    currentEnv,
		"deviceName": devName,
		//"deviceName": config.ClusterName,
	}, &cookie)
	if err != nil {
		return functions.NewE(err)
	}

	if _, err := GetFromResp[bool](respData); err != nil {
		return functions.NewE(err)
	} else {
		return nil
	}
}
