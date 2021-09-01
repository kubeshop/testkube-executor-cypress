package executor

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/kubtest-executor-cypress/internal/pkg/repository/result"
	"github.com/kubeshop/kubtest-executor-cypress/pkg/runner"
	"github.com/kubeshop/kubtest/pkg/worker"

	"github.com/kubeshop/kubtest/pkg/server"

	executorServer "github.com/kubeshop/kubtest/pkg/executor/server"
)

// ConcurrentExecutions per node
const ConcurrentExecutions = 4

// NewExecutor returns new CypressExecutor instance
func NewExecutor(resultRepository result.Repository) Executor {
	var httpConfig server.Config
	envconfig.Process("EXECUTOR", &httpConfig)

	e := Executor{
		HTTPServer: server.NewServer(httpConfig),
		Repository: resultRepository,
		Worker:     worker.NewWorker(resultRepository, &runner.CypressRunner{}),
	}

	return e
}

type Executor struct {
	server.HTTPServer
	Repository result.Repository
	Worker     worker.Worker
}

// Init initialize ExecutorAPI server
func (p *Executor) Init() {

	executions := p.Routes.Group("/executions")

	// add standard start/get handlers from kubtest executor server library
	// they will push and get from worker queue storage
	executions.Post("/", executorServer.StartExecution(p.HTTPServer, p.Repository))
	executions.Get("/:id", executorServer.GetExecution(p.HTTPServer, p.Repository))
}

func (p Executor) Run() error {
	// get executions channel
	executionsQueue := p.Worker.PullExecutions()
	// pass channel to worker
	p.Worker.Run(executionsQueue)

	// run server (blocks process/returns error)
	return p.HTTPServer.Run()
}
