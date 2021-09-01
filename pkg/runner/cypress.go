package runner

import (
	"fmt"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	"github.com/kubeshop/kubtest/pkg/git"
	"github.com/kubeshop/kubtest/pkg/process"
)

func NewCypressRunner() *CypressRunner {
	return &CypressRunner{}
}

// CypressRunner - implements runner interface used in worker to start test execution
type CypressRunner struct {
}

func (r *CypressRunner) Run(execution kubtest.Execution) (result kubtest.ExecutionResult) {

	// make some validation
	err := r.Validate(execution)
	if err != nil {
		return result.Err(err)
	}

	repo := execution.Repository

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

	// TODO move to mapper
	// TODO add result mapping to ExecutionResult
	// map output to Execution result
	return kubtest.ExecutionResult{
		Status:    kubtest.ExecutionStatusSuceess,
		RawOutput: string(out),
	}
}

// Validate checks if Execution has valid data in context of Cypress executor
// Cypress executor runs currently only based on cypress project
func (r *CypressRunner) Validate(execution kubtest.Execution) error {

	if execution.Repository == nil {
		return fmt.Errorf("cypress executor handle only repository based tests, but repository is nil")
	}

	if execution.Repository.Path == "" {
		return fmt.Errorf("can't find repository path in params, repo:%+v", execution.Repository)
	}

	if execution.Repository.Branch == "" {
		return fmt.Errorf("can't find branch in params, repo:%+v", execution.Repository)
	}

	return nil
}
