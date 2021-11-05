package main

import (
	"os"

	"github.com/kubeshop/testkube-executor-cypress/pkg/runner"
	"github.com/kubeshop/testkube/pkg/runner/agent"
)

func main() {
	agent.Run(runner.NewCypressRunner(), os.Args)
}
