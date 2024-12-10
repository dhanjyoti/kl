package fileclient

import fn "github.com/kloudlite/kl/pkg/functions"

var NoEnvSelected = fn.Errorf("no selected environment")

func (fc *fclient) CurrentEnv() (string, error) {
	return "", nil
}

func (fc *fclient) SelectEnv(string) error {
	return nil
}
