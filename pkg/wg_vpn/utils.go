package wg_vpn

import (
	fn "github.com/kloudlite/kl/pkg/functions"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func GenerateWireGuardKeys() (wgtypes.Key, wgtypes.Key, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, fn.Errorf("failed to generate private key: %w", err)
	}
	publicKey := privateKey.PublicKey()

	return privateKey, publicKey, nil
}
