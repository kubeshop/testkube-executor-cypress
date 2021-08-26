package runner

import (
	"io"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
)

// Runner interface to abstract runners implementations
type Runner interface {
	// Run returns output as string (for now probably we could have other needs?)
	Run(input io.Reader, params map[string]string) kubtest.ExecutionResult
}
