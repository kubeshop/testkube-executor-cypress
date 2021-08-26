package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kubeshop/kubtest-executor-template/internal/app/executor"
	"github.com/kubeshop/kubtest-executor-template/internal/pkg/repository/result"
	"github.com/kubeshop/kubtest-executor-template/internal/pkg/storage"
)

const DatabaseName = "template-executor"

type MongoConfig struct {
	DSN string `envconfig:"EXECUTOR_MONGO_DSN" default:"mongodb://localhost:27017"`
	DB  string `envconfig:"EXECUTOR_MONGO_DB" default:"template-executor"`
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

	exec := executor.NewTemplateExecutor(result.NewMongoRespository(db))
	exec.Init()
	panic(exec.Run())

}
