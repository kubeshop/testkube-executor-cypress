package runner

import (
	"io"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
)

// ExampleRunner for template - change me to some valid runner
type ExampleRunner struct {
}

func (r *ExampleRunner) Run(input io.Reader, params map[string]string) kubtest.ExecutionResult {
	return kubtest.ExecutionResult{
		Status:    kubtest.ExecutionStatusSuceess,
		RawOutput: "exmaple test output",
	}
}
