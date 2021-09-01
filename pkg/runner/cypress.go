package runner

import (
	"fmt"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	"github.com/kubeshop/kubtest/pkg/git"
	"github.com/kubeshop/kubtest/pkg/process"
)

// CypressRunner for cypress - change me to some valid runner implements kubtest/pkg/runner.Runner interface
type CypressRunner struct {
}

func (r *CypressRunner) Run(execution kubtest.Execution) (result kubtest.ExecutionResult) {

	repo := execution.Repository

	if repo.Path == "" {
		return result.Err(fmt.Errorf("can't find repository path in params, repo:%+v", repo))
	}

	if repo.Branch == "" {
		return result.Err(fmt.Errorf("can't find branch in params, repo:%+v", repo))
	}

	// checkout repo
	outputDir, err := git.PartialCheckout(repo.Uri, repo.Path, repo.Branch)
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
