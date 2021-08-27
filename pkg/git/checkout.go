package git

import (
	"fmt"
	"io/ioutil"

	"github.com/kubeshop/kubtest/pkg/process"
)

// Partial checkout will checkout only given directory from Git repository
func PartialCheckout(repo, dir string) error {

	tmpDir, err := ioutil.TempDir("", "kubtest-scripts")
	if err != nil {
		return err
	}

	out, err := process.Execute("git",
		"clone",
		"--depth", "1",
		"--filter", "blob:none",
		"--sparse",
		repo, tmpDir)

	fmt.Printf("%s, %s\n", out, err)

	out, err = process.Execute("git", "sparse-checkout", "set", dir)
	fmt.Printf("%s, %s\n", out, err)

	return nil
}
