package runner

import (
	"fmt"
	"net/url"
	"path/filepath"

	junit "github.com/joshdk/go-junit"
	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/git"
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
	GitUsername     string // RUNNER_GITUSERNAME
	GitToken        string // RUNNER_GITTOKEN
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
	uri := repo.Uri
	if r.Params.GitUsername != "" && r.Params.GitToken != "" {
		gitURI, err := url.Parse(uri)
		if err != nil {
			return result, err
		}

		gitURI.User = url.UserPassword(r.Params.GitUsername, r.Params.GitToken)
		uri = gitURI.String()
	}

	// checkout repo
	outputDir, err := git.PartialCheckout(uri, repo.Path, repo.Branch)
	if err != nil {
		return result, err
	}

	// be gentle to different cypress versions, run from local npm deps
	_, err = executor.Run(outputDir, "npm", "install")
	if err != nil {
		return result, err
	}

	junitReportPath := filepath.Join(outputDir, "results/junit.xml")
	args := []string{"run", "--reporter", "junit", "--reporter-options", fmt.Sprintf("mochaFile=%s,toConsole=false", junitReportPath)}
	// append args from execution
	args = append(args, execution.Args...)

	// run cypress inside repo directory ignore execution error in case of failed test
	out, err := executor.Run(outputDir, "./node_modules/cypress/bin/cypress", args...)
	suites, serr := junit.IngestFile(junitReportPath)
	result = MapJunitToExecutionResults(out, suites)

	if r.Params.ScrapperEnabled {
		fmt.Println("Scrapper enabled fetching videos and snapshots")
		client := minio.NewClient(r.Params.Endpoint, r.Params.AccessKeyID, r.Params.SecretAccessKey, r.Params.Location, r.Params.Token, r.Params.Ssl) // create storage client
		err := client.Connect()
		if err != nil {
			// TODO fix this one - should log or maybe introduce some warning status for test results
			fmt.Println("error occured creating minio client", err) // maybe we should consider the run failed since it is not able to save artefacts
		}

		directories := []string{
			filepath.Join(outputDir, "cypress/videos"),
			filepath.Join(outputDir, "cypress/screenshots"),
		}

		err = client.ScrapeArtefacts(execution.Id, directories...)
		if err != nil {
			// TODO fix this one
			fmt.Println("error occured while scrapping artefacts", err) // maybe we should consider the run failed since it is not able to save artefacts
		}
	}

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
