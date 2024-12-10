package shell

import (
	"os"
	"path"
	"strings"
)

func mount(mounts map[string]string, mountpath string) error {

	for k, v := range mounts {
		key := strings.TrimPrefix(k, "$kl_mounts")

		prefixedPath := path.Join(mountpath, key)

		// fn.Logf(text.Yellow("mounting %s at %s\n"), key, prefixedPath)
		os.MkdirAll(path.Dir(prefixedPath), 0o700)

		if err := os.WriteFile(prefixedPath, []byte(v), 0o700); err != nil {
			return err
		}
	}

	return nil
}
