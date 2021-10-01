package main

import (
	"fmt"

	"github.com/kubeshop/kubtest-executor-cypress/pkg/runner"
	"github.com/kubeshop/kubtest/pkg/api/v1/kubtest"
)

func main() {
	runner := runner.NewCypressRunner()
	repoURI := "https://github.com/kubeshop/kubtest-executor-cypress.git"
	result := runner.Run(kubtest.Execution{
		Params: map[string]string{"testparam": "testvalue"},
		Repository: &kubtest.Repository{
			Uri:    repoURI,
			Branch: "jacek/feature/json-output",
			Path:   "examples",
		},
	})

	fmt.Printf("%+v\n", result)

}
