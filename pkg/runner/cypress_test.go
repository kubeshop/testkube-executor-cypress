package runner

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {

	runner := CypressRunner{}
	repoInput := bytes.NewReader([]byte("https://github.com/cirosantilli/test-git-partial-clone-big-small"))
	result := runner.Run(repoInput, map[string]string{
		"dir":    "examples",
		"branch": "jacek/feature/git-checkout",
	})

	fmt.Printf("%+v\n", result)

	t.Fail()

}
