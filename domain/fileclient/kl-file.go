package fileclient

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	confighandler "github.com/kloudlite/kl/pkg/config-handler"
	"github.com/kloudlite/kl/pkg/functions"
	fn "github.com/kloudlite/kl/pkg/functions"
)

type KLFileType struct {
	Version    string `json:"version" yaml:"version"`
	DefaultEnv string `json:"defaultEnv,omitempty" yaml:"defaultEnv,omitempty"`
	TeamName   string `json:"teamName,omitempty" yaml:"teamName,omitempty"`

	Packages  []string `json:"packages" yaml:"packages"`
	Libraries []string `json:"libraries" yaml:"libraries"`

	EnvVars EnvVars `json:"envVars" yaml:"envVars"`
	Mounts  Mounts  `json:"mounts" yaml:"mounts"`
	Ports   []int   `json:"ports" yaml:"ports"`

	// packagesMap  map[string]int `json:"-"`
	// librariesMap map[string]int `json:"-"`
	// ConfigFile   string         `json:"-"`
}

type HashData map[string]string

type Lockfile struct {
	Packages  HashData `json:"packages" yaml:"packages"`
	Libraries HashData `json:"libraries" yaml:"libraries"`
}

func (k *Lockfile) Save() error {
	if k == nil {
		return fmt.Errorf("lockfile is nil")
	}

	if err := confighandler.WriteConfig(fmt.Sprintf("%s.lock", getConfigPath()), *k, 0o644); err != nil {
		fn.PrintError(err)
		return functions.NewE(err)
	}

	return nil
}

func (c *fclient) GetLockfile() (*Lockfile, error) {
	filePath := getConfigPath()

	kllockfile, err := confighandler.ReadConfig[Lockfile](fmt.Sprintf("%s.lock", filePath))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		return &Lockfile{
			Packages:  HashData{},
			Libraries: HashData{},
		}, nil
	}

	return kllockfile, nil
}

func (k *KLFileType) Save() error {
	if k == nil {
		return fmt.Errorf("klfile is nil")
	}

	if err := confighandler.WriteConfig(getConfigPath(), *k, 0o644); err != nil {
		fn.PrintError(err)
		return functions.NewE(err)
	}

	return nil

}

const (
	defaultKLFile = "kl.yml"
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

	if dir == "/" {
		return "", fmt.Errorf("config file not found")
	}

	return lookupKLFile(filepath.Dir(dir))
}

func getConfigPath() string {
	file, ok := os.LookupEnv("KLCONFIG_PATH")
	if !ok {
		wd, _ := os.Getwd()
		var err error
		file, err = lookupKLFile(wd)
		if err != nil {
			return defaultKLFile
		}
	}

	return file
}

func (c *fclient) WriteKLFile(fileObj KLFileType) error {
	if err := confighandler.WriteConfig(getConfigPath(), fileObj, 0o644); err != nil {
		fn.PrintError(err)
		return functions.NewE(err)
	}

	return nil
}

func (c *fclient) GetKlFile() (*KLFileType, error) {
	klfile, err := c.getKlFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fn.Errorf("no kl.yml found, please run `kl init` to initialize kl.yml")
		}

		return nil, functions.NewE(err)
	}
	return klfile, nil
}

func (c *fclient) getKlFile() (*KLFileType, error) {
	filePath := getConfigPath()

	klfile, err := confighandler.ReadConfig[KLFileType](filePath)
	if err != nil {
		return nil, functions.NewE(err, "failed to read klfile")
	}

	return klfile, nil
}
