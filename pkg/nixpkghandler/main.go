package nixpkghandler

import (
	"github.com/kloudlite/kl/domain/fileclient"
	"github.com/spf13/cobra"
)

type PackageClient interface {
	AddPackage(name, hash string) error
	RemovePackage(name string) error

	AddLibrary(name, hash string) error
	RemoveLibrary(name string) error

	// used for listing search results
	Search(name string) (*SearchResults, error)

	// used for fzf search
	Find(pname string) (string, string, error)
	SyncLockfile() error
}

type pkgHandler struct {
	cmd *cobra.Command
	fc  fileclient.FileClient
}

func New(cmd *cobra.Command) (PackageClient, error) {
	fc, err := fileclient.New()
	if err != nil {
		return nil, err
	}

	return &pkgHandler{
		cmd: cmd,
		fc:  fc,
	}, nil
}
