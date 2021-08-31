package runner

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	"github.com/kubeshop/kubtest/pkg/git"
	"github.com/kubeshop/kubtest/pkg/process"
)

// CypressRunner for cypress - change me to some valid runner
type CypressRunner struct {
}

func (r *CypressRunner) Run(input io.Reader, params map[string]string) (result kubtest.ExecutionResult) {
	repoBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return result.Err(err)
	}
	repo := string(repoBytes)

	dir, ok := params["dir"]
	if !ok {
		return result.Err(fmt.Errorf("can't find directory ('dir') in params"))
	}

	branch, ok := params["branch"]
	if !ok {
		return result.Err(fmt.Errorf("can't find directory ('dir') in params"))
	}

	// checkout repo
	outputDir, err := git.PartialCheckout(repo, dir, branch)
	if err != nil {
		return result.Err(err)
	}

	// be gentle to different cypress versions, run from local npm deps
	_, err = process.ExecuteInDir(outputDir, "npm", "install")
	if err != nil {
		return result.Err(err)
	}

	// run cypress inside repo directory
	out, err := process.ExecuteInDir(outputDir, "./node_modules/cypress/bin/cypress", "run")
	if err != nil {
		return result.Err(err)
	}

	// map output to Execution result
	// TODO move to mapper
	return kubtest.ExecutionResult{
		Status:    kubtest.ExecutionStatusSuceess,
		RawOutput: string(out),
	}
}
