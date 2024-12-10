package fileclient

import (
	"bytes"
	"os"
	"path"

	confighandler "github.com/kloudlite/kl/pkg/config-handler"
	fn "github.com/kloudlite/kl/pkg/functions"
	"github.com/kloudlite/kl/pkg/ui/text"
)

type CacheKLConfig struct {
	Hash    string
	EnvVars map[string]string
	Mounts  map[string]string
}

type wsContextData struct {
	EnvName string `json:"envName"`

	Cache *CacheKLConfig `json:"cache"`
}

func (w wsContext) GetCache() *CacheKLConfig {
	return w.Cache
}

func (w wsContext) SetCache(cache *CacheKLConfig) error {
	w.Cache = cache
	return w.handler.Write()
}

type WsContext interface {
	SetEnv(string) error
	GetEnv() (string, error)
	GetCache() *CacheKLConfig
	SetCache(cache *CacheKLConfig) error
}

func (w wsContext) GetEnv() (string, error) {
	if w.EnvName == "" {
		return "", fn.Errorf("env not found")
	}

	s, err := getCtxData()
	if err != nil {
		return "", err
	}

	menv, err := s.GetEnv()
	if err != nil {
		return "", err
	}

	if menv != w.EnvName {
		return "", fn.Errorf("selected env %s is not same as current working directory env %s, please change selected env using %s", text.Yellow(menv), text.Yellow(w.EnvName), text.Blue("kl use env"))
	}

	return w.EnvName, nil
}

func (w wsContext) SetEnv(env string) error {

	s, err := getCtxData()
	if err != nil {
		return err
	}

	if err := s.SetEnv(env); err != nil {
		return err
	}

	w.EnvName = env
	return w.handler.Write()
}

type wsContext struct {
	*wsContextData
	handler confighandler.Config[wsContextData]
}

func getNewWsContext() (WsContext, error) {
	cpath, err := assertConfigPath()
	if err != nil {
		return nil, err
	}

	cdir := path.Dir(cpath)

	cachePath := path.Join(cdir, ".kl")
	if err := os.MkdirAll(cachePath, os.ModePerm); err != nil {
		return nil, err
	}

	// ensure .kl in .gitignore
	b, err := os.ReadFile(path.Join(cdir, ".gitignore"))
	if err != nil {
		b = []byte{}
	}

	if !bytes.Contains(b, []byte(".kl")) {
		b = append(b, []byte("\n.kl\n")...)
		if err := os.WriteFile(path.Join(cdir, ".gitignore"), b, 0o644); err != nil {
			return nil, err
		}
	}

	chandler := confighandler.GetHandler[wsContextData](path.Join(cachePath, "config.yaml"))

	cdata, _ := chandler.Read()

	ch := &wsContext{
		handler:       chandler,
		wsContextData: cdata,
	}

	return ch, nil
}
