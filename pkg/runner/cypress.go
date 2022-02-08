package runner

import (
	"fmt"
	"path/filepath"
	"strings"

	junit "github.com/joshdk/go-junit"
	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/output"
	"github.com/kubeshop/testkube/pkg/executor/scrapper"
)

type Params struct {
	Endpoint        string // RUNNER_ENDPOINT
	AccessKeyID     string // RUNNER_ACCESSKEYID
	SecretAccessKey string // RUNNER_SECRETACCESSKEY
	Location        string // RUNNER_LOCATION
	Token           string // RUNNER_TOKEN
	Ssl             bool   // RUNNER_SSL
	ScrapperEnabled bool   // RUNNER_SCRAPPERENABLED
	GitUsername     string // RUNNER_GITUSERNAME
	GitToken        string // RUNNER_GITTOKEN
}

func NewCypressRunner() *CypressRunner {
	var params Params
	err := envconfig.Process("runner", &params)
	if err != nil {
		panic(err.Error())
	}

	runner := &CypressRunner{
		Fetcher: content.NewFetcher(),
		Scrapper: scrapper.NewScrapper(
			params.Endpoint,
			params.AccessKeyID,
			params.SecretAccessKey,
			params.Location,
			params.Token,
			params.Ssl,
		),
		Params: params,
	}

	return runner
}

// CypressRunner - implements runner interface used in worker to start test execution
type CypressRunner struct {
	Params   Params
	Fetcher  content.ContentFetcher
	Scrapper *scrapper.Scrapper
}

func (r *CypressRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {

	// make some validation
	err = r.Validate(execution)
	if err != nil {
		return result, err
	}

	path, err := r.Fetcher.Fetch(execution.Content)
	if err != nil {
		return result, err
	}

	if execution.Content.IsFile() {
		output.PrintEvent("using file", execution)

		// TODO add cypress project structure
		// TODO checkout this repo with `skeleton` path
		// TODO overwrite skeleton/cypress/integration/test.js
		//      file with execution content git file
		return result, fmt.Errorf("passing cypress script as single file not implemented yet")
	}

	// be gentle to different cypress versions, run from local npm deps
	_, err = executor.Run(path, "npm", "install")
	if err != nil {
		return result, fmt.Errorf("npm install error: %w", err)
	}

	envVars := make([]string, 0, len(execution.Params))
	for key, value := range execution.Params {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	junitReportPath := filepath.Join(path, "results/junit.xml")
	args := []string{"run", "--reporter", "junit", "--reporter-options", fmt.Sprintf("mochaFile=%s,toConsole=false", junitReportPath),
		"--env", strings.Join(envVars, ",")}

	// append args from execution
	args = append(args, execution.Args...)

	// run cypress inside repo directory ignore execution error in case of failed test
	out, err := executor.Run(path, "./node_modules/cypress/bin/cypress", args...)
	suites, serr := junit.IngestFile(junitReportPath)
	result = MapJunitToExecutionResults(out, suites)

	// scrape artifacts first even if there are errors above
	if r.Params.ScrapperEnabled {
		directories := []string{
			filepath.Join(path, "cypress/videos"),
			filepath.Join(path, "cypress/screenshots"),
		}
		err := r.Scrapper.Scrape(execution.Id, directories)
		if err != nil {
			return result.WithErrors(fmt.Errorf("scrape artifacts error: %w", err)), nil
		}
	}

	return result.WithErrors(err, serr), nil
}

// Validate checks if Execution has valid data in context of Cypress executor
// Cypress executor runs currently only based on cypress project
func (r *CypressRunner) Validate(execution testkube.Execution) error {

	if execution.Content == nil {
		return fmt.Errorf("can't find any content to run in execution data: %+v", execution)
	}

	if execution.Content.Repository == nil {
		return fmt.Errorf("cypress executor handle only repository based tests, but repository is nil")
	}

	if execution.Content.Repository.Path == "" {
		return fmt.Errorf("can't find repository path in params, repo:%+v", execution.Content.Repository)
	}

	if execution.Content.Repository.Branch == "" {
		return fmt.Errorf("can't find branch in params, repo:%+v", execution.Content.Repository)
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
