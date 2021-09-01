package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/kubtest-executor-cypress/pkg/runner"
	"github.com/kubeshop/kubtest/pkg/executor/repository/result"
	"github.com/kubeshop/kubtest/pkg/executor/repository/storage"
	"github.com/kubeshop/kubtest/pkg/executor/server"
	"github.com/kubeshop/kubtest/pkg/ui"
)

const DatabaseName = "cypress-executor"

type MongoConfig struct {
	DSN        string `envconfig:"EXECUTOR_MONGO_DSN" default:"mongodb://localhost:27017"`
	DB         string `envconfig:"EXECUTOR_MONGO_DB" default:"cypress-executor"`
	Collection string `envconfig:"EXECUTOR_MONGO_COLLECTION" default:"executions"`
}

var cfg MongoConfig

func init() {
	envconfig.Process("mongo", &cfg)
}

func main() {
	db, err := storage.GetMongoDataBase(cfg.DSN, cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := result.NewMongoRespository(db, cfg.Collection)
	runner := runner.NewCypressRunner()
	exec := server.NewExecutor(repo, runner)

	ui.ExitOnError("Running executor", exec.Init().Run())
}
