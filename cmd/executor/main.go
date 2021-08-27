package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/kubtest-executor-cypress/internal/app/executor"
	"github.com/kubeshop/kubtest-executor-cypress/internal/pkg/repository/result"
	"github.com/kubeshop/kubtest-executor-cypress/internal/pkg/storage"
)

const DatabaseName = "cypress-executor"

type MongoConfig struct {
	DSN string `envconfig:"EXECUTOR_MONGO_DSN" default:"mongodb://localhost:27017"`
	DB  string `envconfig:"EXECUTOR_MONGO_DB" default:"cypress-executor"`
}

var Config MongoConfig

func init() {
	envconfig.Process("mongo", &Config)
}

func main() {

	db, err := storage.GetMongoDataBase(Config.DSN, Config.DB)
	if err != nil {
		panic(err)
	}

	exec := executor.NewExecutor(result.NewMongoRespository(db))
	exec.Init()
	panic(exec.Run())

}
