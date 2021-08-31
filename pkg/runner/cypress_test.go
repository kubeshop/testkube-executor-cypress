package runner

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestRun(t *testing.T) {

	// Can't run it in my default install
	os.Setenv("CYPRESS_CACHE_FOLDER", "/Users/exu/tmp")

	runner := CypressRunner{}
	repoURI := "https://github.com/kubeshop/kubtest-executor-cypress.git"
	repoInput := bytes.NewReader([]byte(repoURI))
	result := runner.Run(repoInput, map[string]string{
		"dir":    "examples",
		"branch": "jacek/feature/git-checkout",
	})

	fmt.Printf("%+v\n", result)

	t.Fail()

}
