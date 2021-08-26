package executor

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	"github.com/stretchr/testify/assert"
)

func TestTemplateExecutor_StartExecution(t *testing.T) {

	t.Run("runs newman runner command", func(t *testing.T) {
		// given
		executor := GetTestExecutor(t)

		req := httptest.NewRequest(
			"POST",
			"/v1/executions/",
			strings.NewReader(`{"type": "template/collection", "metadata": "{\"info\":{\"name\":\"kubtestExampleCollection\"}}"}`),
		)

		// when
		resp, err := executor.Mux.Test(req)
		assert.NoError(t, err)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

}

func GetTestExecutor(t *testing.T) TemplateExecutor {
	templateExecutor := NewTemplateExecutor(&RepoMock{
		Object: kubtest.Execution{Id: "1"},
	})
	templateExecutor.Init()

	return templateExecutor
}

// r RepoMock
type RepoMock struct {
	Object kubtest.Execution
	Error  error
}

func (r *RepoMock) Get(ctx context.Context, id string) (result kubtest.Execution, err error) {
	return r.Object, r.Error
}

func (r *RepoMock) Insert(ctx context.Context, result kubtest.Execution) (err error) {
	return r.Error
}

func (r *RepoMock) QueuePull(ctx context.Context) (result kubtest.Execution, err error) {
	return r.Object, r.Error
}

func (r *RepoMock) Update(ctx context.Context, result kubtest.Execution) (err error) {
	return r.Error
}
