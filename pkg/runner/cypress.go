package runner

import (
	"io"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
)

// CypressRunner for template - change me to some valid runner
type CypressRunner struct {
}

func (r *CypressRunner) Run(input io.Reader, params map[string]string) kubtest.ExecutionResult {
	return kubtest.ExecutionResult{
		Status:    kubtest.ExecutionStatusSuceess,
		RawOutput: "exmaple test output",
	}
}
