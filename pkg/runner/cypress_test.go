package runner

import (
	"fmt"
	"os"
	"testing"

	"github.com/kubeshop/kubtest/pkg/api/v1/kubtest"
)

func TestRun(t *testing.T) {
	// t.Skip("move this test to e2e test suite with valid environment setup")

	// Can't run it in my default install on mac
	os.Setenv("CYPRESS_CACHE_FOLDER", os.TempDir())

	runner := CypressRunner{}
	repoURI := "https://github.com/kubeshop/kubtest-executor-cypress.git"
	result := runner.Run(kubtest.Execution{
		Repository: &kubtest.Repository{
			Uri:    repoURI,
			Branch: "main",
			Path:   "examples",
		},
	})

	fmt.Printf("%+v\n", result)

	t.Fail()

}
