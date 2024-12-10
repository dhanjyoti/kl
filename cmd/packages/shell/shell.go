package shell

import (
	"fmt"
	"os"
	"path"

	"github.com/kloudlite/kl/domain/apiclient"
	"github.com/kloudlite/kl/domain/fileclient"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/nixpkghandler"
	"github.com/kloudlite/kl/pkg/ui/text"
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
	kshflag, ok := os.LookupEnv("KL_SHELL")
	if ok && kshflag == "true" {
		return fmt.Errorf(text.Red("You are already in an active kl shell.\nRun `exit` before calling `kl shell` again. Shell inception is not supported."))
	}

	apic, err := apiclient.New()
	if err != nil {
		return err
	}

	fc, err := fileclient.New()
	if err != nil {
		return err
	}

	kpath, err := fc.GetConfigPath()
	if err != nil {
		return err
	}

	mountpath := path.Join(path.Dir(kpath), ".kl", "mounts")

	pc, err := nixpkghandler.New(cmd)
	if err != nil {
		return err
	}

	var envMap, mountMap map[string]string

	ck, err := getCache()
	if err == nil {
		envMap = ck.EnvVars
		mountMap = ck.Mounts
	} else {
		fn.Log(text.Yellow("cache not found, refetching"))
		envMap, mountMap, err = apic.GetLoadMaps()
		if err != nil {
			fn.Warn(err)
			envMap = make(map[string]string)
			mountMap = make(map[string]string)
		}

		if err := setCache(&fileclient.CacheKLConfig{
			Mounts:  mountMap,
			EnvVars: envMap,
		}); err != nil {
			return err
		}
	}

	if err := mount(mountMap, mountpath); err != nil {
		return err
	}

	if err := pc.SyncLockfile(); err != nil {
		return err
	}

	lockFile, err := apic.GetFClient().GetLockfile()

	pkgs := make([]string, 0, len(lockFile.Packages))
	for _, v := range lockFile.Packages {
		pkgs = append(pkgs, fmt.Sprintf("nixpkgs/%s", v))
	}

	libs := make([]string, 0, len(lockFile.Libraries))
	for _, v := range lockFile.Libraries {
		libs = append(libs, fmt.Sprintf("nixpkgs/%s", v))
	}

	envs := make([]string, 0, len(envMap))
	for k, v := range envMap {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	if err := NixShell(cmd.Context(), ShellArgs{
		Shell:     os.Getenv("SHELL"),
		EnvVars:   append(envs, "KL_SHELL=true", fmt.Sprintf("kl_mounts=%s", mountpath)),
		Packages:  pkgs,
		Libraries: libs,
	}); err != nil {
		return err
	}

	return nil
}
