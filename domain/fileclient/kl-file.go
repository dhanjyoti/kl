package fileclient

import (
	"errors"
	"fmt"
	"os"
	"path"
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

	cpath, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := confighandler.WriteConfig(fmt.Sprintf("%s.lock", cpath), *k, 0o644); err != nil {
		fn.PrintError(err)
		return functions.NewE(err)
	}

	return nil
}

func (c *fclient) GetLockfile() (*Lockfile, error) {
	cpath, err := assertConfigPath()
	if err != nil {
		return nil, err
	}

	kllockfile, err := confighandler.ReadConfig[Lockfile](fmt.Sprintf("%s.lock", cpath))
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

	cpath, err := getConfigPath()
	if err != nil && err.Error() == "config file not found" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		cpath = path.Join(wd, defaultKLFile)
	}

	if err := confighandler.WriteConfig(cpath, *k, 0o644); err != nil {
		return functions.NewE(err)
	}

	return nil

}

const (
	defaultKLFile = "kl.yaml"
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

func (fc *fclient) GetConfigPath() (string, error) {
	return assertConfigPath()
}

func GetKlPath() (string, error) {
	return assertConfigPath()
}

func assertConfigPath() (string, error) {
	cpath, err := getConfigPath()
	if err != nil {
		return "", fn.Errorf("can't find kl.yaml file in the working directory, please run `kl init` to initialize kl.yaml")
	}

	if _, err := os.Stat(cpath); os.IsNotExist(err) {
		return "", err
	}

	return cpath, nil
}

func getConfigPath() (string, error) {
	file, ok := os.LookupEnv("KLCONFIG_PATH")
	if !ok {
		wd, _ := os.Getwd()
		var err error

		file, err = lookupKLFile(wd)
		if err != nil {
			return "", err
		}
	}

	return file, nil
}

func (c *fclient) WriteKLFile(fileObj KLFileType) error {
	cpath, err := assertConfigPath()
	if err != nil {
		return err
	}

	if cpath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		cpath = filepath.Join(wd, defaultKLFile)
	}

	if err := confighandler.WriteConfig(cpath, fileObj, 0o644); err != nil {
		fn.PrintError(err)
		return functions.NewE(err)
	}

	return nil
}

func (c *fclient) GetKlFile() (*KLFileType, error) {
	klfile, err := getKlFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fn.Errorf("no kl.yml found, please run `kl init` to initialize kl.yml")
		}

		return nil, functions.NewE(err)
	}
	return klfile, nil
}

func getKlFile() (*KLFileType, error) {
	filePath, err := assertConfigPath()
	if err != nil {
		return nil, err
	}

	klfile, err := confighandler.ReadConfig[KLFileType](filePath)
	if err != nil {
		return nil, functions.NewE(err, "failed to read klfile")
	}

	return klfile, nil
}
