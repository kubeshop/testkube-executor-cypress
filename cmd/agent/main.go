package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kubeshop/testkube-executor-cypress/pkg/runner"
	"github.com/kubeshop/testkube/pkg/executor/agent"
	"github.com/kubeshop/testkube/pkg/executor/output"
)

var (
	dependecy = flag.String("dependency", "npm", "package manager")
)

func main() {
	flag.Parse()

	r, err := runner.NewCypressRunner(*dependecy)
	if err != nil {
		output.PrintError(os.Stderr, fmt.Errorf("could not initialize runner: %w", err))
		os.Exit(1)
	}
	agent.Run(r, os.Args)
}
