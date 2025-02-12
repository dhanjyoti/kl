package add

import (
	"bufio"
	"fmt"
	"github.com/kloudlite/kl/cmd/box/boxpkg"
	"github.com/kloudlite/kl/cmd/box/boxpkg/hashctrl"
	"github.com/kloudlite/kl/domain/apiclient"
	"github.com/kloudlite/kl/domain/fileclient"
	"github.com/kloudlite/kl/pkg/functions"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/fzf"
	"github.com/kloudlite/kl/pkg/ui/spinner"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var mountCommand = &cobra.Command{
	Use:   "config-mount [path]",
	Short: "add file mount to your kl-config file by selection from the all the [ config | secret ] available in current environemnt",
	Long: `
	This command will help you to add mounts to your kl-config file.
	You can add a config or secret to your kl-config file by providing the path of the config/secret you want to add.
	`,
	Example: `
  kl add config-mount [path] --config=<config_name>	# add mount from config.
  kl add config-mount [path] --secret=<secret_name>	# add secret from secret.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fc, err := fileclient.New()
		if err != nil {
			fn.PrintError(err)
			return
		}

		apic, err := apiclient.New()
		if err != nil {
			fn.PrintError(err)
			return
		}

		filePath := fn.ParseKlFile(cmd)

		if filePath == "" {
			filePath = "/home/kl/workspace/kl.yml"
		}

		klFile, err := fc.GetKlFile(filePath)
		if err != nil {
			fn.PrintError(err)
			return
		}

		err = selectConfigMount(apic, fc, *klFile, cmd, args)
		if err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func selectConfigMount(apic apiclient.ApiClient, fc fileclient.FileClient, klFile fileclient.KLFileType, cmd *cobra.Command, args []string) error {

	//TODO: add changes to the klbox-hash file
	c := cmd.Flag("config").Value.String()
	s := cmd.Flag("secret").Value.String()

	var cOrs fileclient.CSType
	cOrs = ""

	if c != "" || s != "" {

		if c != "" {
			cOrs = fileclient.ConfigType
		} else {
			cOrs = fileclient.SecretType
		}

	} else {
		csName := []fileclient.CSType{fileclient.ConfigType, fileclient.SecretType}
		cOrsValue, err := fzf.FindOne(
			csName,
			//func(i int) string {
			//	return csName[i]
			//},
			func(item fileclient.CSType) string {
				return string(item)
			},
			fzf.WithPrompt("Mount from Config/Secret >"),
		)
		if err != nil {
			return fn.NewE(err)
		}

		cOrs = fileclient.CSType(*cOrsValue)
	}

	items := make([]apiclient.ConfigORSecret, 0)
	if cOrs == fileclient.ConfigType {
		currentTeam, err := fc.CurrentTeamName()
		if err != nil {
			return err
		}
		currentEnv, err := apic.EnsureEnv()
		if err != nil {
			fn.PrintError(err)
			return err
		}
		configs, e := apic.ListConfigs(currentTeam, currentEnv.Name)

		if e != nil {
			return e
		}

		for _, c := range configs {
			items = append(items, apiclient.ConfigORSecret{
				Entries: c.Data,
				Name:    c.Metadata.Name,
			})
		}

	} else {
		currentTeam, err := fc.CurrentTeamName()
		if err != nil {
			return err
		}
		currentEnv, err := apic.EnsureEnv()
		if err != nil {
			fn.PrintError(err)
			return err
		}
		secrets, e := apic.ListSecrets(currentTeam, currentEnv.Name)

		if e != nil {
			return e
		}

		for _, c := range secrets {
			items = append(items, apiclient.ConfigORSecret{
				Entries: c.StringData,
				Name:    c.Metadata.Name,
			})
		}
	}

	if len(items) == 0 {
		return fn.Errorf("no %ss created yet on server ", cOrs)
	}

	selectedItem := apiclient.ConfigORSecret{}

	if c != "" || s != "" {
		csId := func() string {
			if c != "" {
				return c
			}
			return s
		}()

		for _, co := range items {
			if co.Name == csId {
				selectedItem = co
				break
			}
		}

		return fn.Errorf("provided %s name not found", cOrs)
	} else {
		selectedItemVal, err := fzf.FindOne(
			items,
			func(item apiclient.ConfigORSecret) string {
				return item.Name
			},
			fzf.WithPrompt(fmt.Sprintf("Select %s >", cOrs)),
		)

		if err != nil {
			fn.PrintError(err)
		}

		selectedItem = *selectedItemVal
	}

	key, err := fzf.FindOne(func() []string {
		res := make([]string, 0)
		for k := range selectedItem.Entries {
			res = append(res, k)
		}
		return res
	}(), func(item string) string {
		return item
	}, fzf.WithPrompt("Select Config/Secret >"))

	if err != nil {
		return fn.NewE(err)
	}

	spinner.Client.Pause()
	fn.Printf("path of the config file, (eg: /tmp/sample): ")
	path, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fn.PrintError(err)
	}
	path = strings.TrimSpace(path)
	defer spinner.Client.Resume()

	matchedIndex := -1
	for i, fe := range klFile.Mounts {
		if fe.Path == path {
			matchedIndex = i
		}
	}

	fe := klFile.Mounts.GetMounts()

	if matchedIndex == -1 {
		fe = append(fe, fileclient.FileEntry{
			Type: cOrs,
			Path: path,
			Name: selectedItem.Name,
			Key:  *key,
		})
	} else {
		fe[matchedIndex] = fileclient.FileEntry{
			Type: cOrs,
			Path: path,
			Name: selectedItem.Name,
			Key:  *key,
		}
	}

	klFile.Mounts.AddMounts(fe)
	if err := fc.WriteKLFile(klFile); err != nil {
		return fn.NewE(err)
	}

	fn.Log("added mount to your kl-file")

	wpath, err := os.Getwd()
	if err != nil {
		return fn.NewE(err)
	}

	if err = hashctrl.SyncBoxHash(apic, fc, wpath); err != nil {
		return fn.NewE(err)
	}

	cl, err := boxpkg.NewClient(cmd, nil)
	if err != nil {
		return functions.NewE(err)
	}

	if err := cl.ConfirmBoxRestart(); err != nil {
		return functions.NewE(err)
	}

	return nil
}

func init() {
	mountCommand.Flags().StringP("config", "", "", "config name")
	mountCommand.Flags().StringP("secret", "", "", "secret name")
	fn.WithKlFile(mountCommand)
}
