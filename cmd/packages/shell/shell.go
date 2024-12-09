package shell

import (
	"fmt"
	"os"

	"github.com/kloudlite/kl/domain/apiclient"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/nixpkghandler"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "shell",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Shell(cmd, args); err != nil {
			fn.PrintError(err)
		}
	},
}

func Shell(cmd *cobra.Command, args []string) error {

	apic, err := apiclient.New()
	if err != nil {
		return err
	}

	pc, err := nixpkghandler.New(cmd)
	if err != nil {
		return err
	}

	envMap, mountMap, err := apic.GetLoadMaps()
	if err != nil {
		fn.Warn(err)
		envMap = make(map[string]string)
		mountMap = make(map[string]string)
	}

	fmt.Println(mountMap)

	if err := pc.SyncLockfile(); err != nil {
		return err
	}

	lockFile, err := apic.GetFileClient().GetLockfile()

	pkgs := make([]string, 0)
	for _, v := range lockFile.Packages {
		pkgs = append(pkgs, fmt.Sprintf("nixpkgs/%s", v))
	}

	libs := make([]string, 0)
	for _, v := range lockFile.Libraries {
		libs = append(libs, fmt.Sprintf("nixpkgs/%s", v))
	}

	envs := make([]string, 0)
	for k, v := range envMap {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	fmt.Println(envs)

	if err := NixShell(cmd.Context(), ShellArgs{
		Shell:     os.Getenv("SHELL"),
		EnvVars:   append(envs, "KL_SHELL=true"),
		Packages:  pkgs,
		Libraries: libs,
	}); err != nil {
		return err
	}

	return nil
}
