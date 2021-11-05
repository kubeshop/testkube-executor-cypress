package runner

import (
	"fmt"
	"os"
	"path/filepath"

	junit "github.com/joshdk/go-junit"
	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/git"
	"github.com/kubeshop/testkube/pkg/process"
	"github.com/kubeshop/testkube/pkg/runner/output"
	"github.com/kubeshop/testkube/pkg/storage/minio"
)

type Params struct {
	Endpoint        string // RUNNER_ENDPOINT
	AccessKeyID     string // RUNNER_ACCESSKEYID
	SecretAccessKey string // RUNNER_SECRETACCESSKEY
	Location        string // RUNNER_LOCATION
	Token           string // RUNNER_TOKEN
	Ssl             bool   // RUNNER_SSL
	ScrapperEnabled bool   // RUNNER_SCRAPPERENABLED
}

func NewCypressRunner() *CypressRunner {
	runner := &CypressRunner{}

	err := envconfig.Process("runner", &runner.Params)
	if err != nil {
		panic(err.Error())
	}

	return runner
}

// CypressRunner - implements runner interface used in worker to start test execution
type CypressRunner struct {
	Params Params
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

	// wrap stdout lines into JSON chunks we want it to have common interface for agent
	// stdin <- testkube.Execution, stdout <- stream of json logs
	// LoggedExecuteInDir will put wrapped JSON output to stdout AND get RAW output into out var
	// json logs can be processed later on watch of pod logs
	writer := output.NewJSONWrapWriter(os.Stdout)

	// be gentle to different cypress versions, run from local npm deps
	_, err = process.LoggedExecuteInDir(outputDir, writer, "npm", "install")
	if err != nil {
		return result, err
	}

	junitReportPath := filepath.Join(outputDir, "results/junit.xml")
	args := []string{"run", "--reporter", "junit", "--reporter-options", fmt.Sprintf("mochaFile=%s,toConsole=false", junitReportPath)}
	for k, v := range execution.Params {
		args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
	}

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

	if r.Params.ScrapperEnabled {
		client, err := minio.NewClient(r.Params.Endpoint, r.Params.AccessKeyID, r.Params.SecretAccessKey, r.Params.Location, r.Params.Token, r.Params.Ssl) // create storage client
		if err != nil {
			fmt.Println("error occured creating minio client") // maybe we should consider the run failed since it is not able to save artefacts
		}

		err = client.ScrapeArtefacts(execution.Id, "cypress/")
		if err != nil {
			fmt.Println("error occured while scrapping artefacts") // maybe we should consider the run failed since it is not able to save artefacts
		}
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
