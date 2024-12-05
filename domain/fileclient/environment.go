package fileclient

import (
	"os"

	fn "github.com/kloudlite/kl/pkg/functions"
)

var NoEnvSelected = fn.Errorf("no selected environment")

func (f *fclient) SelectEnv(ev Env) error {
	k, err := GetExtraData()
	if err != nil {
		return fn.NewE(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return fn.NewE(err)
	}

	if k.SelectedEnvs == nil {
		k.SelectedEnvs = map[string]*Env{}
	}

	k.SelectedEnvs[dir] = &ev

	return SaveExtraData(k)
}

func (f *fclient) SelectEnvOnPath(pth string, ev Env) error {
	k, err := GetExtraData()
	if err != nil {
		return fn.NewE(err)
	}

	if k.SelectedEnvs == nil {
		k.SelectedEnvs = map[string]*Env{}
	}

	k.SelectedEnvs[pth] = &ev

	return SaveExtraData(k)
}

func (f *fclient) EnvOfPath(pth string) (*Env, error) {
	c, err := GetExtraData()
	if err != nil {
		return nil, fn.NewE(err)
	}

	if c.SelectedEnvs == nil || c.SelectedEnvs[pth] == nil {
		return nil, fn.NewE(NoEnvSelected)
	}

	return c.SelectedEnvs[pth], nil
}

func (f *fclient) CurrentEnv() (*Env, error) {
	c, err := GetExtraData()
	if err != nil {
		return nil, fn.NewE(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, fn.NewE(err)
	}

	if c.SelectedEnvs == nil {
		return nil, fn.NewE(NoEnvSelected)
	}

	if c.SelectedEnvs[dir] == nil {
		return nil, fn.NewE(NoEnvSelected)
	}

	return c.SelectedEnvs[dir], nil
}
