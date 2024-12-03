package shell

import (
	"crypto/md5"
	"fmt"
	"slices"
)

func GenHash(pkl *ParsedKLConfig) string {
	h := md5.New()
	slices.Sort(pkl.EnvVars)
	for i := range pkl.EnvVars {
		h.Write([]byte(pkl.EnvVars[i]))
	}

	mk := make([]string, 0, len(pkl.Mounts))
	for k := range pkl.Mounts {
		mk = append(mk, k)
	}

	for _, k := range mk {
		h.Write(pkl.Mounts[k])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
