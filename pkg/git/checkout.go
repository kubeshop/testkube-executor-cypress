package git

import (
	"fmt"
	"io/ioutil"

	"github.com/kubeshop/kubtest/pkg/process"
)

// TODO consider checkout by some go-git library to limit external deps in docker container
// Partial checkout will checkout only given directory from Git repository
func PartialCheckout(repo, dir string) (outputDir string, err error) {

	tmpDir, err := ioutil.TempDir("", "kubtest-scripts")
	if err != nil {
		return tmpDir, err
	}

	fmt.Printf("%+v\n", tmpDir)

	out, err := process.ExecuteInDir(tmpDir, "git",
		"clone",
		"--depth", "1",
		"--filter", "blob:none",
		"--sparse",
		repo, "repo")

	fmt.Printf("%s, %s\n", out, err)

	out, err = process.ExecuteInDir(tmpDir+"/repo", "git", "sparse-checkout", "set", dir)
	fmt.Printf("%s, %s\n", out, err)

	return tmpDir + "/repo/" + dir, nil
}
