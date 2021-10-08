package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kubeshop/testkube-executor-cypress/pkg/runner"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("missing input argument")
		os.Exit(1)
	}

	script := args[1]

	e := testkube.Execution{}
	json.Unmarshal([]byte(script), &e)
	runner := runner.NewCypressRunner()
	result := runner.Run(e)
	fmt.Println(result)
	fmt.Printf("$$$%s$$$", e.Id)
}
