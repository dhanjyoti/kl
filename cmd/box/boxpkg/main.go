package boxpkg

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kloudlite/kl/k3s"

	"github.com/kloudlite/kl/domain/apiclient"
	"github.com/kloudlite/kl/domain/fileclient"
	"github.com/kloudlite/kl/flags"

	dockerclient "github.com/docker/docker/client"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/spf13/cobra"
)

type client struct {
	cli        *dockerclient.Client
	cmd        *cobra.Command
	args       []string
	foreground bool
	verbose    bool
	cwd        string

	containerName string

	env *fileclient.Env

	fc     fileclient.FileClient
	apic   apiclient.ApiClient
	klfile *fileclient.KLFileType
	k3s    k3s.K3sClient
}

type BoxClient interface {
	SyncProxy(config ProxyConfig) error
	Stop() error
	Restart() error
	Start() error
	Ssh() error
	Reload() error
	PrintBoxes([]Cntr) error
	ListAllBoxes() ([]Cntr, error)
	Info() error
	Exec([]string, io.Writer) error

	ConfirmBoxRestart() error
	StartWgContainer() error
	StopContainer() error
}

func (c *client) Context() context.Context {
	return c.cmd.Context()
}

func NewClient(cmd *cobra.Command, args []string) (BoxClient, error) {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fn.NewE(err)
	}

	fc, err := fileclient.New()
	if err != nil {
		return nil, fn.NewE(err)
	}

	apic, err := apiclient.New()
	if err != nil {
		return nil, fn.NewE(err)
	}

	foreground := fn.ParseBoolFlag(cmd, "foreground")
	cwd, _ := os.Getwd()

	hash := md5.New()
	hash.Write([]byte(cwd))
	contName := fmt.Sprintf("klbox-%s", fmt.Sprintf("%x", hash.Sum(nil))[:8])

	klFile, err := fc.GetKlFile("")
	if err != nil {
		return nil, fn.NewE(err)
	}

	k3sClient, err := k3s.NewClient(cmd)
	if err != nil {
		return nil, fn.NewE(err)
	}

	env, err := fc.EnvOfPath(cwd)
	if err != nil && errors.Is(err, fileclient.NoEnvSelected) {
		env := &fileclient.Env{
			SSHPort: 0,
		}
		if klFile.DefaultEnv != "" && klFile.TeamName != "" {
			environment, err := apic.GetEnvironment(klFile.TeamName, klFile.DefaultEnv)
			if err != nil {
				return nil, fn.NewE(err)
			}
			env.Name = environment.Metadata.Name
		}

		data, err := fileclient.GetExtraData()
		if err != nil {
			return nil, fn.NewE(err)
		}
		if data.SelectedEnvs == nil {
			data.SelectedEnvs = map[string]*fileclient.Env{
				cwd: env,
			}
		} else {
			data.SelectedEnvs[cwd] = env
		}
		if err := fileclient.SaveExtraData(data); err != nil {
			return nil, fn.NewE(err)
		}
	}

	return &client{
		cli:           cli,
		cmd:           cmd,
		args:          args,
		foreground:    foreground,
		verbose:       flags.IsVerbose,
		cwd:           cwd,
		containerName: contName,
		env:           env,
		fc:            fc,
		apic:          apic,
		klfile:        klFile,
		k3s:           k3sClient,
	}, nil
}
