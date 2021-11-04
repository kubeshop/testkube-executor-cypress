package runner

import (
	"fmt"
	"os"
	"path/filepath"

	junit "github.com/joshdk/go-junit"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/git"
	"github.com/kubeshop/testkube/pkg/process"
	"github.com/kubeshop/testkube/pkg/runner/output"
)

func NewCypressRunner() *CypressRunner {
	return &CypressRunner{}
}

// CypressRunner - implements runner interface used in worker to start test execution
type CypressRunner struct {
}

func (r *CypressRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {

	// make some validation
	err = r.Validate(execution)
	if err != nil {
		return result, err
	}

	repo := execution.Repository

	// checkout repo
	outputDir, err := git.PartialCheckout(repo.Uri, repo.Path, repo.Branch)
	if err != nil {
		return result, err
	}

	// be gentle to different cypress versions, run from local npm deps
	_, err = process.LoggedExecuteInDir(outputDir, os.Stdout, "npm", "install")
	if err != nil {
		return result, err
	}

	junitReportPath := filepath.Join(outputDir, "results/junit.xml")
	args := []string{"run", "--reporter", "junit", "--reporter-options", fmt.Sprintf("mochaFile=%s,toConsole=false", junitReportPath)}
	for k, v := range execution.Params {
		args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
	}

	// wrap stdout lines into JSON chunks we want it to have common interface for agent
	// stdin <- testkube.Execution, stdout <- stream of json logs
	// LoggedExecuteInDir will put wrapped JSON output to stdout AND get RAW output into out var
	// json logs can be processed later on watch of pod logs
	writer := output.NewJSONWrapWriter(os.Stdout)

	// run cypress inside repo directory ignore execution error in case of failed test
	out, err := process.LoggedExecuteInDir(outputDir, writer, "./node_modules/cypress/bin/cypress", args...)
	suites, serr := junit.IngestFile(junitReportPath)
	result = MapJunitToExecutionResults(out, suites)

	// handle errors if any
	if err != nil {
		return result.Err(err), nil
	}
	if serr != nil {
		return result.Err(serr), nil
	}

	return
}

// Validate checks if Execution has valid data in context of Cypress executor
// Cypress executor runs currently only based on cypress project
func (r *CypressRunner) Validate(execution testkube.Execution) error {

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

func MapJunitToExecutionResults(out []byte, suites []junit.Suite) (result testkube.ExecutionResult) {
	status := testkube.SUCCESS_ExecutionStatus
	result.Status = &status
	result.Output = string(out)
	result.OutputType = "text/plain"

	for _, suite := range suites {
		for _, test := range suite.Tests {

			result.Steps = append(
				result.Steps,
				testkube.ExecutionStepResult{
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
		return string(testkube.SUCCESS_ExecutionStatus)
	default:
		return string(testkube.ERROR__ExecutionStatus)
	}
}
