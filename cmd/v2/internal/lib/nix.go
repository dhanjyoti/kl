package lib

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InstallPackage(interactive bool, pkgs ...string) ([]string, error) {
	c := exec.Command("sh", "-c", fmt.Sprintf("nix shell %s --command echo downloaded", strings.Join(pkgs, " ")))

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		return nil, err
	}

	return nil, nil
}

type ShellArgs struct {
	EnvVars   []string
	Packages  []string
	Libraries []string
}

func pathExists(p string) error {
	_, err := os.Stat(p)
	if err != nil {
		return err
	}
	return nil
}

func NixShell(ctx context.Context, args ShellArgs) error {
	ev := append(os.Environ(), args.EnvVars...)

	libPaths := make([]string, 0, len(args.Libraries))
	var includes []string

	for _, lib := range args.Libraries {
		// nix eval $package --raw
		// nix-store --query --references $package
		c := exec.CommandContext(ctx, "nix", "eval", lib, "--raw")
		b, err := c.CombinedOutput()
		if err != nil {
			return err
		}

		// ev = append(ev, fmt.Sprintf("LD_LIBRARY_PATH=%s:%s", string(b)+"/lib", os.Getenv("LD_LIBRARY_PATH")))
		if pathExists(string(b)+"/include") != nil {
			includes = append(includes, string(b)+"/include")
		}

		c2 := exec.CommandContext(ctx, "nix-store", "--query", "--references", string(b))
		b2, err := c2.CombinedOutput()
		if err != nil {
			return err
		}
		lines := strings.Split(string(b2), "\n")
		// fmt.Printf("b2: %v %d\n", lines, len(lines))

		libPaths = append(libPaths, string(b)+"/lib")
		for i := range lines {
			if len(strings.TrimSpace(lines[i])) > 0 {
				if pathExists(lines[i]+"/lib") != nil {
					libPaths = append(libPaths, lines[i]+"/lib")
				}
			}
		}
	}

	libPaths = createSet(libPaths)
	includes = createSet(includes)

	fmt.Printf("LD_LIBRARY_PATH: %s\n#######\n", strings.Join(libPaths, ":"))
	ev = append(ev, fmt.Sprintf("LD_LIBRARY_PATH=%s:%s", strings.Join(libPaths, ":"), os.Getenv("LD_LIBRARY_PATH")))
	ev = append(ev, fmt.Sprintf("CPATH=%s:%s", strings.Join(includes, ":"), os.Getenv("CPATH")))
	c := exec.Command("sh", "-c", fmt.Sprintf("nix shell %s", strings.Join(args.Packages, " ")))

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Env = ev

	return c.Run()
}
