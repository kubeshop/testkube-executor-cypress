package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckOut(t *testing.T) {

	repo := "https://github.com/cirosantilli/test-git-partial-clone-big-small"

	err := PartialCheckout(repo, "small")
	assert.NoError(t, err)
	t.Fail()
}
