package initp

import (
	"os"
	"path/filepath"

	"github.com/kloudlite/kl/domain/fileclient"
	fn "github.com/kloudlite/kl/pkg/functions"

	"github.com/kloudlite/kl/pkg/ui/text"

	"github.com/spf13/cobra"
)

var InitCommand = &cobra.Command{
	Use:   "init",
	Short: "initialize a kl-config file",
	Long:  `use this command to initialize a kl-config file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleInit(); err != nil {
			fn.PrintError(err)
			return
		}
	},
}

func handleInit() error {
	fc, err := fileclient.New()
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	files := []string{"kl.yml", "kl.yaml"}

	configPath := ""
	for _, name := range files {
		fi, err := os.Stat(filepath.Join(cwd, name))
		if err != nil {
			continue
		}
		if fi.IsDir() {
			continue
		}
		configPath = filepath.Join(cwd, name)
		break
	}

	if configPath != "" {
		fn.Printf(text.Yellow("workspace is already initilized. Do you want to override? [y/N]: "))
		if !fn.Confirm("Y", "N") {
			return fn.Errorf("file initialization aborted")
		}
	}

	newKlFile := fileclient.KLFileType{
		Version:  "v1",
		Packages: []string{"neovim", "git"},
	}

	if err := fc.WriteKLFile(newKlFile); err != nil {
		fn.PrintError(err)
	} else {
		fn.Printf(text.Green("workspace initialized successfully.\n"))
	}

	return nil
}

func init() {
}
