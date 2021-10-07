package runner

import (
	"fmt"
	"os"
	"path/filepath"

	junit "github.com/joshdk/go-junit"
	"github.com/kubeshop/kubtest/pkg/api/v1/kubtest"
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
	_, err = process.LoggedExecuteInDir(outputDir, os.Stdout, "npm", "install")
	if err != nil {
		return result.Err(err)
	}

	junitReportPath := filepath.Join(outputDir, "results/junit.xml")
	args := []string{"run", "--reporter", "junit", "--reporter-options", fmt.Sprintf("mochaFile=%s,toConsole=false", junitReportPath)}
	for k, v := range execution.Params {
		args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
	}

	// run cypress inside repo directory
	out, err := process.LoggedExecuteInDir(outputDir, os.Stdout, "./node_modules/cypress/bin/cypress", args...)
	if err != nil {
		return result.Err(err)
	}

	suites, err := junit.IngestFile(junitReportPath)
	if err != nil {
		return result.Err(err)
	}

	return MapJunitToExecutionResults(out, suites)
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

func MapJunitToExecutionResults(out []byte, suites []junit.Suite) (result kubtest.ExecutionResult) {
	status := kubtest.SUCCESS_ExecutionStatus
	result.Status = &status
	result.Output = string(out)
	result.OutputType = "text/plain"

	for _, suite := range suites {
		for _, test := range suite.Tests {

			result.Steps = append(
				result.Steps,
				kubtest.ExecutionStepResult{
					Name:     fmt.Sprintf("%s - %s", suite.Name, test.Name),
					Duration: test.Duration.String(),
					Status:   MapStatus(test.Status),
				})
		}

		// TODO parse sub suites recursively

	}

	return result
}

func MapStatus(in junit.Status) (out string) {
	switch string(in) {
	case "passed":
		return string(kubtest.SUCCESS_ExecutionStatus)
	default:
		return string(kubtest.ERROR__ExecutionStatus)
	}
}
