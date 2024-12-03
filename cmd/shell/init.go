package shell

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kloudlite/kl/domain/fileclient"
	confighandler "github.com/kloudlite/kl/pkg/config-handler"
	"github.com/kloudlite/kl/pkg/functions"
	"github.com/spf13/cobra"
)

func lookupKLFile(dir string) (string, error) {
	files := []string{"kl.yml", "kl.yaml"}
	for _, name := range files {
		fi, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			continue
			// return nil, err
		}
		if fi.IsDir() {
			continue
			// return nil, fmt.Errorf("config file: %s is a directory, must be a file", fi.Name())
		}
		return filepath.Join(dir, name), nil
	}

	return lookupKLFile(filepath.Dir(dir))
}

func findKLFile() (string, error) {
	file, ok := os.LookupEnv("KLCONFIG_PATH")
	if !ok {
		wd, _ := os.Getwd()
		var err error
		file, err = lookupKLFile(wd)
		if err != nil {
			return "", err
		}
	}

	slog.Debug("got", "file", file)
	return file, nil
}

type ParsedKLConfig struct {
	ConfigFile string
	CacheDir   string

	EnvVars []string
	Mounts  map[string][]byte

	Hash string
}

func (pkl *ParsedKLConfig) WriteHash() error {
	h := GenHash(pkl)
	pkl.Hash = h
	return os.WriteFile(filepath.Join(pkl.CacheDir, "hash.txt"), []byte(h), 0o5000)
}

func (pkl *ParsedKLConfig) CompareHash() bool {
	b, err := os.ReadFile(filepath.Join(pkl.CacheDir, "hash.txt"))
	if err != nil {
		return false
	}

	return pkl.Hash == string(b)
}

// INFO: read more at cmd/box/boxpkg/hashctrl/main.go:231 (generatePersistedEnv)
func parseKLConfig(file string) (*ParsedKLConfig, error) {
	cfg, err := confighandler.ReadConfig[fileclient.KLFileType](file)
	if err != nil {
		return nil, functions.NewE(err, "failed to read klfile")
	}

	p := &ParsedKLConfig{
		ConfigFile: file,
		CacheDir:   filepath.Join(filepath.Dir(file), ".kl"),
		EnvVars:    make([]string, 0, len(cfg.EnvVars)),
	}

	// ensures kl config directory
	if err := os.Mkdir(p.CacheDir, 0o700); err != nil {
		if !os.IsExist(err) {
			panic("failed to create .kl directory")
		}
	}

	// TODO: add this cachedir to `.gitignore`

	envVars, mounts, err := ParseEnvVarsAndMounts(cfg)
	if err != nil {
		panic(err)
	}
	p.EnvVars = append(p.EnvVars, envVars...)
	p.Mounts = mounts

	for _, v := range cfg.EnvVars {
		if v.Value != nil {
			p.EnvVars = append(p.EnvVars, fmt.Sprintf("%s=%s", v.Key, *v.Value))
		}
	}

	if err := p.WriteHash(); err != nil {
		return nil, err
	}

	return p, nil
}

/*
TODO: hash of parsed EnvVars and Mounts
*/

var Command = &cobra.Command{
	Use:   "shell",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		klf, err := findKLFile()
		if err != nil {
			panic(err)
		}

		pkl, err := parseKLConfig(klf)
		if err != nil {
			panic(err)
		}

		for k, v := range pkl.Mounts {
			rk := filepath.Join(pkl.CacheDir, "mounts", k)
			_, err := os.Stat(rk)
			if err != nil {
				slog.Error("got", "err", err, "condition", errors.Is(err, os.ErrNotExist))
				if !errors.Is(err, os.ErrNotExist) {
					panic(err)
				}

				if err := os.MkdirAll(filepath.Dir(rk), 0o760); err != nil {
					if !os.IsExist(err) {
						panic(fmt.Errorf("failed to create mount dir (%s)", filepath.Dir(rk)))
					}
				}

				if err := os.WriteFile(rk, v, 0o444); err != nil {
					panic(err)
				}
				continue
			}
			// panic(fmt.Errorf("filepath: %s, already exists", fi.Name()))
		}

		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		env := make([]string, 0, len(os.Environ())+len(pkl.EnvVars)+1)
		env = append(env, os.Environ()...)
		env = append(env, pkl.EnvVars...)
		env = append(env, "KL_SHELL=true", "KL_HASH="+pkl.Hash)

		c := exec.Command(shell)
		c.Env = env
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Printf("Failed to start shell process: %v\n", err)
		}
	},
}
